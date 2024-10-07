package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/machship/graphql"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			schema, err := graphql.NewSchema(graphql.SchemaConfig{
				Query: graphql.NewObject(graphql.ObjectConfig{
					Name: "RootQuery",
					Fields: graphql.Fields{
						"hello": {
							Type: graphql.String,
							Resolve: func(p graphql.ResolveParams) (any, error) {
								return "world", nil
							},
						},
					},
				}),
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create schema: %s", err)
				return
			}
			runtime.KeepAlive(schema)
		}()
	}

	wg.Wait()
}
