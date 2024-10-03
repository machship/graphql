package graphql_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRace(t *testing.T) {
	tempdir := filepath.Join(os.TempDir(), "race")
	if err := os.MkdirAll(tempdir, 0755); err != nil {
		t.Fatalf("Failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tempdir)

	filename := filepath.Join(tempdir, "example.go")
	err := os.WriteFile(filename, []byte(`
		package main

		import (
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
					schema, _ := graphql.NewSchema(graphql.SchemaConfig{
						Query: graphql.NewObject(graphql.ObjectConfig{
							Name: "RootQuery",
							Fields: graphql.Fields{
								"hello": &graphql.Field{
									Type: graphql.String,
									Resolve: func(p graphql.ResolveParams) (any, error) {
										return "world", nil
									},
								},
							},
						}),
					})
					runtime.KeepAlive(schema)
				}()
			}

			wg.Wait()
		} 
	`), 0755)
	if err != nil {
		t.Fatal(err)
	}

	result, err := exec.Command("go", "run", "-race", filename).CombinedOutput()
	if err != nil || len(result) != 0 {
		t.Log(string(result))
		t.Fatal(err)
	}
}
