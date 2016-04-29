package filesystem

import (
	"os"
	"testing"

	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestDeploymentDelete(t *testing.T) {
	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("proj", "fancy")

	deployment := Deployment(tmp)
	deployment.Remove("proj", "fancy")

	dir := tmp.Path("proj", "fancy")
	if isDirectory(dir) {
		t.Logf("`Delete` should have deleted '%s', but didn't.", dir)
		t.Fail()
	}
}

func TestDeploymentDeleteSkipMaster(t *testing.T) {
	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("proj", "master")

	deployment := Deployment(tmp)
	deployment.Remove("proj", "master")

	dir := tmp.Path("proj", "master")
	if !isDirectory(dir) {
		t.Logf("`Delete` shouldn't have deleted '%s', but did.", dir)
		t.Fail()
	}
}

func TestDeploymentDeleteSkip(t *testing.T) {
	for _, sub := range []string{"current", "release", "shared"} {
		tmp := testutil.CreateDirectory(t)
		defer tmp.Clean()

		tmp.Branch("proj", "fancy")

		os.Remove(tmp.Path("proj", "fancy", sub))
		deployment := Deployment(tmp)
		deployment.Remove("proj", "fancy")

		dir := tmp.Path("proj", "fancy")
		if !isDirectory(dir) {
			t.Logf("`Delete` shouldn't have deleted '%s', but did.", dir)
			t.Fail()
		}
	}
}
