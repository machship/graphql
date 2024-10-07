package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/machship/graphql"
	"github.com/machship/graphql/gqlerrors"
	"github.com/machship/graphql/language/ast"
	"github.com/machship/graphql/language/parser"
)

func TestParse(t *testing.T, query string) *ast.Document {
	astDoc, err := parser.Parse(parser.ParseParams{
		Source: query,
		Options: parser.ParseOptions{
			// include source, for error reporting
			NoSource: false,
		},
	})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	return astDoc
}
func TestExecute(t *testing.T, ep graphql.ExecuteParams) *graphql.Result {
	return graphql.Execute(ep)
}

func Diff(want, got any) []string {
	return []string{fmt.Sprintf("\ngot:\n%v", got), fmt.Sprintf("\nwant:\n%v\n", want)}
}

func ASTToJSON(t *testing.T, a ast.Node) any {
	b, err := json.Marshal(a)
	if err != nil {
		t.Fatalf("Failed to marshal Node %v", err)
	}
	var f any
	err = json.Unmarshal(b, &f)
	if err != nil {
		t.Fatalf("Failed to unmarshal Node %v", err)
	}
	return f
}

func ContainSubsetSlice(super []any, sub []any) bool {
	if len(sub) == 0 {
		return true
	}
subLoop:
	for _, subVal := range sub {
		found := false
	innerLoop:
		for _, superVal := range super {
			if subVal, ok := subVal.(map[string]any); ok {
				if superVal, ok := superVal.(map[string]any); ok {
					if ContainSubset(superVal, subVal) {
						found = true
						break innerLoop
					} else {
						continue
					}
				} else {
					return false
				}

			}
			if subVal, ok := subVal.([]any); ok {
				if superVal, ok := superVal.([]any); ok {
					if ContainSubsetSlice(superVal, subVal) {
						found = true
						break innerLoop
					} else {
						continue
					}
				} else {
					return false
				}
			}
			if reflect.DeepEqual(superVal, subVal) {
				found = true
				break innerLoop
			}
		}
		if !found {
			return false
		}
		continue subLoop
	}
	return true
}

func ContainSubset(super map[string]any, sub map[string]any) bool {
	if len(sub) == 0 {
		return true
	}
	for subKey, subVal := range sub {
		if superVal, ok := super[subKey]; ok {
			switch superVal := superVal.(type) {
			case []any:
				if subVal, ok := subVal.([]any); ok {
					if !ContainSubsetSlice(superVal, subVal) {
						return false
					}
				} else {
					return false
				}
			case map[string]any:
				if subVal, ok := subVal.(map[string]any); ok {
					if !ContainSubset(superVal, subVal) {
						return false
					}
				} else {
					return false
				}
			default:
				if !reflect.DeepEqual(superVal, subVal) {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func EqualErrorMessage(expected, result *graphql.Result, i int) bool {
	return expected.Errors[i].Message == result.Errors[i].Message
}

func EqualFormattedError(exp, act gqlerrors.FormattedError) bool {
	if exp.Message != act.Message {
		return false
	}
	if !reflect.DeepEqual(exp.Locations, act.Locations) {
		return false
	}
	if !reflect.DeepEqual(exp.Path, act.Path) {
		return false
	}
	if !reflect.DeepEqual(exp.Extensions, act.Extensions) {
		return false
	}
	return true
}

func EqualFormattedErrors(expected, actual []gqlerrors.FormattedError) bool {
	if len(expected) != len(actual) {
		return false
	}
	for i := range expected {
		if !EqualFormattedError(expected[i], actual[i]) {
			return false
		}
	}
	return true
}

func EqualResults(expected, result *graphql.Result) bool {
	if !reflect.DeepEqual(expected.Data, result.Data) {
		return false
	}
	return EqualFormattedErrors(expected.Errors, result.Errors)
}

// testCtxKey is a type to avoid key collisions in context values.
type testCtxKey string

// ContextWithValue inserts value into the context.
func ContextWithValue(ctx context.Context, k string, v any) context.Context {
	return context.WithValue(ctx, testCtxKey(k), v)
}

// ContextValue retrieves a value from the context.
func ContextValue(ctx context.Context, k string) any {
	return ctx.Value(testCtxKey(k))
}

// RootDir returns an absolute path to the root of the module. It fatally fails
// the test if the runtime information could not be obtained.
func RootDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to retrieve runtime information")
	}

	return filepath.Join(filepath.Dir(filename), "..")
}

// PathFromRoot calls [RootDir] and creates a path with the result and the provided
// slice of path components.
func PathFromRoot(t *testing.T, components ...string) string {
	t.Helper()
	components = append([]string{RootDir(t)}, components...)
	return filepath.Join(components...)
}
