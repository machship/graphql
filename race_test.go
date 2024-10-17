package graphql_test

import (
	"os/exec"
	"testing"

	"github.com/machship/graphql/testutil"
)

func TestRace(t *testing.T) {
	path := testutil.PathFromRoot(t, "testutil", "race")
	result, err := exec.Command("go", "run", "-race", path).CombinedOutput()

	if err != nil || len(result) != 0 {
		t.Logf("%s", result)
		t.Fatal(err)
	}
}
