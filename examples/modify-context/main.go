package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/machship/graphql"
	"github.com/machship/graphql/testutil"
)

type User struct {
	ID int `json:"id"`
}

var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (any, error) {
				rootValue := p.Info.RootValue.(map[string]any)
				if rootValue["data-from-parent"] == "ok" &&
					rootValue["data-before-execution"] == "ok" {
					user := p.Source.(User)
					return user.ID, nil
				}
				return nil, nil
			},
		},
	},
})

func main() {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"users": &graphql.Field{
					Type: graphql.NewList(UserType),
					Resolve: func(p graphql.ResolveParams) (any, error) {
						rootValue := p.Info.RootValue.(map[string]any)
						rootValue["data-from-parent"] = "ok"
						result := []User{
							{ID: 1},
						}
						return result, nil

					},
				},
			},
		}),
	})
	if err != nil {
		log.Fatal(err)
	}
	ctx := testutil.ContextWithValue(context.Background(), "currentUser", User{ID: 100})
	// Instead of trying to modify context within a resolve function, use:
	// `graphql.Params.RootObject` is a mutable optional variable and available on
	// each resolve function via: `graphql.ResolveParams.Info.RootValue`.
	rootObject := map[string]any{
		"data-before-execution": "ok",
	}
	result := graphql.Do(graphql.Params{
		Context:       ctx,
		RequestString: "{ users { id } }",
		RootObject:    rootObject,
		Schema:        schema,
	})
	b, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(b)) // {"data":{"users":[{"id":1}]}}
}
