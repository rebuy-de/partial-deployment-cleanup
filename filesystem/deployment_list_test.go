package filesystem

import (
	"reflect"
	"testing"

	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestList(t *testing.T) {
	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("proj", "fancy")
	tmp.Branch("proj", "master")
	tmp.Branch("proj", "foo")
	tmp.Branch("proj", "bar")

	deployment := Deployment(tmp)
	branches, err := deployment.GetBranches("proj")
	if err != nil {
		t.Fatal(err.Error())
	}

	expect := []string{"bar", "fancy", "foo", "master"}

	if !reflect.DeepEqual(branches, expect) {
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Obtained: %#v", branches)
	}
}
