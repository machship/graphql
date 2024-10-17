package graphql_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/machship/graphql"
	"github.com/machship/graphql/testutil"
)

func TestAppliedDirectives(t *testing.T) {
	lengthDirective := graphql.NewDirective(graphql.DirectiveConfig{
		Name:        "length",
		Description: "Used to specify the minimum and/or maximum length for an input field or argument.",
		Locations: []string{
			graphql.DirectiveLocationFieldDefinition,
		},
		Args: graphql.FieldConfigArgument{
			"min": {
				Type:        graphql.Int,
				Description: "If specified, specifies the minimum length that the input field or argument must have.",
			},
			"max": {
				Type:        graphql.Int,
				Description: "If specified, specifies the maximum length that the input field or argument must have.",
			},
		},
	})
	directives := []*graphql.Directive{
		lengthDirective,
	}

	type Droid struct {
		ID          int64
		Name        string
		CustomField string
	}

	r2d2 := Droid{
		ID:   1,
		Name: "R2-D2",
	}
	root := graphql.NewObject(graphql.ObjectConfig{
		Name: "root",
		Fields: graphql.Fields{
			"hero": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "DroidType",
					Fields: graphql.Fields{
						"id": {
							Type: graphql.ID,
							Resolve: func(p graphql.ResolveParams) (any, error) {
								return r2d2.ID, nil
							},
						},
						"name": {
							Type: graphql.String,
							Resolve: func(p graphql.ResolveParams) (any, error) {
								return r2d2.Name, nil
							},
						},
						"customField": {
							Type: graphql.String,
							Resolve: func(p graphql.ResolveParams) (any, error) {
								return r2d2.CustomField, nil
							},
							Directives: []*graphql.AppliedDirective{
								lengthDirective.Apply([]*graphql.DirectiveArgument{
									{
										Name:  "min",
										Value: 103,
									},
									{
										Name:  "max",
										Value: 999,
									},
								}),
							},
						},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (any, error) {
					return []any{
						r2d2,
					}, nil
				},
			},
		},
	})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:      root,
		Directives: directives,
	})
	if err != nil {
		t.Fatalf("Failed to create new schema: %s", err)
	}

	out := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: testutil.AppliedDirectiveIntrospectionQuery,
		// RequestString: in,
		Context: context.Background(),
	})

	if out.HasErrors() {
		t.Fatalf("Failed to execute request: %v", out.Errors)
	}

	b, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		t.Fatalf("Failed to marshal result: %s", err)
	}

	t.Logf("\n\n%s\n\n", b)
}
