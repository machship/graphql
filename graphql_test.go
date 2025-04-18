package graphql_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/machship/graphql"
	"github.com/machship/graphql/testutil"
	"github.com/machship/graphql/testutil/starwars"
)

type T struct {
	name      string
	Query     string
	Schema    graphql.Schema
	Expected  any
	Variables map[string]any
}

var Tests = []T{}

func init() {
	Tests = []T{
		{
			Query: `
				query HeroNameQuery {
					hero {
						name
					}
				}
			`,
			Schema: starwars.Schema,
			Expected: &graphql.Result{
				Data: map[string]any{
					"hero": map[string]any{
						"name": "R2-D2",
					},
				},
			},
		},
		{
			name: "xyz",
			Query: `
				query HeroNameAndFriendsQuery {
					hero {
						id
						name
						friends {
							name
						}
					}
				}
			`,
			Schema: starwars.Schema,
			Expected: &graphql.Result{
				Data: map[string]any{
					"hero": map[string]any{
						"id":   "2001",
						"name": "R2-D2",
						"friends": []any{
							map[string]any{
								"name": "Luke Skywalker",
							},
							map[string]any{
								"name": "Han Solo",
							},
							map[string]any{
								"name": "Leia Organa",
							},
						},
					},
				},
			},
		},
		{
			Query: `
				query HumanByIdQuery($id: String!) {
					human(id: $id) {
						name
					}
				}
			`,
			Schema: starwars.Schema,
			Expected: &graphql.Result{
				Data: map[string]any{
					"human": map[string]any{
						"name": "Darth Vader",
					},
				},
			},
			Variables: map[string]any{
				"id": "1001",
			},
		},
	}
}

func TestQuery(t *testing.T) {
	for _, test := range Tests {
		params := graphql.Params{
			Schema:         test.Schema,
			RequestString:  test.Query,
			VariableValues: test.Variables,
		}
		testGraphql(test, params, t)
	}
}

func testGraphql(test T, p graphql.Params, t *testing.T) {
	result := graphql.Do(p)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result, test.Expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", test.Query, testutil.Diff(test.Expected, result))
	}
}

func TestBasicGraphQLExample(t *testing.T) {
	// taken from `graphql-js` README

	helloFieldResolved := func(p graphql.ResolveParams) (any, error) {
		return "world", nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQueryType",
			Fields: graphql.Fields{
				"hello": &graphql.Field{
					Description: "Returns `world`",
					Type:        graphql.String,
					Resolve:     helloFieldResolved,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := "{ hello }"
	var expected any = map[string]any{
		"hello": "world",
	}

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}

func TestThreadsContextFromParamsThrough(t *testing.T) {
	extractFieldFromContextFn := func(p graphql.ResolveParams) (any, error) {
		return testutil.ContextValue(p.Context, p.Args["key"].(string)), nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"value": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: extractFieldFromContextFn,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ value(key:"a") }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       testutil.ContextWithValue(context.Background(), "a", "xyz"),
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]any{"value": "xyz"}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}

}

func TestNewErrorChecksNilNodes(t *testing.T) {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"graphql_is": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (any, error) {
						return "", nil
					},
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected errors: %v", err.Error())
	}
	query := `{graphql_is:great(sort:ByPopularity)}{stars}`
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) == 0 {
		t.Fatalf("expected errors, got: %v", result)
	}
}

func TestEmptyStringIsNotNull(t *testing.T) {
	checkForEmptyString := func(p graphql.ResolveParams) (any, error) {
		arg := p.Args["arg"]
		if arg == nil || arg.(string) != "" {
			t.Errorf("Expected empty string for input arg, got %#v", arg)
		}
		return "yay", nil
	}
	returnEmptyString := func(p graphql.ResolveParams) (any, error) {
		return "", nil
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"checkEmptyArg": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"arg": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: checkForEmptyString,
				},
				"checkEmptyResult": &graphql.Field{
					Type:    graphql.String,
					Resolve: returnEmptyString,
				},
			},
		}),
	})
	if err != nil {
		t.Fatalf("wrong result, unexpected errors: %v", err.Error())
	}
	query := `{ checkEmptyArg(arg:"") checkEmptyResult }`

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	expected := map[string]any{"checkEmptyArg": "yay", "checkEmptyResult": ""}
	if !reflect.DeepEqual(result.Data, expected) {
		t.Errorf("wrong result, query: %v, graphql result diff: %v", query, testutil.Diff(expected, result))
	}
}
