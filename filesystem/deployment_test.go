package filesystem

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func testHelperDeploymentDirectory(t *testing.T, project string, branches ...string) (string, func()) {
	tmp, err := ioutil.TempDir("", "pdtest")
	if err != nil {
		t.Logf(err.Error())
		t.FailNow()
	}

	for _, branch := range branches {
		os.MkdirAll(path.Join(tmp, project, branch), 0755)

		for _, sub := range []string{"current", "release", "shared"} {
			os.Mkdir(path.Join(tmp, project, branch, sub), 0755)
		}
	}

	return tmp, func() {
		os.RemoveAll(tmp)
	}
}

func TestDeploymentDelete(t *testing.T) {
	tmp, fn := testHelperDeploymentDirectory(t, "proj", "fancy")
	defer fn()

	deployment := Deployment(tmp)
	deployment.Delete("proj", "fancy")

	dir := path.Join(tmp, "proj", "fancy")
	if isDirectory(dir) {
		t.Logf("`Delete` should have deleted '%s', but didn't.", dir)
		t.Fail()
	}
}

func TestDeploymentDeleteSkipMaster(t *testing.T) {
	tmp, fn := testHelperDeploymentDirectory(t, "proj", "master")
	defer fn()

	deployment := Deployment(tmp)
	deployment.Delete("proj", "master")

	dir := path.Join(tmp, "proj", "master")
	if !isDirectory(dir) {
		t.Logf("`Delete` shouldn't have deleted '%s', but did.", dir)
		t.Fail()
	}
}

func TestDeploymentDeleteSkip(t *testing.T) {
	for _, sub := range []string{"current", "release", "shared"} {
		tmp, fn := testHelperDeploymentDirectory(t, "proj", "fancy")
		defer fn()

		os.Remove(path.Join(tmp, "proj", "fancy", sub))
		deployment := Deployment(tmp)
		deployment.Delete("proj", "fancy")

		dir := path.Join(tmp, "proj", "fancy")
		if !isDirectory(dir) {
			t.Logf("`Delete` shouldn't have deleted '%s', but did.", dir)
			t.Fail()
		}
	}
}

func TestList(t *testing.T) {
	tmp, def := testHelperDeploymentDirectory(t, "proj", "fancy", "master", "foo", "bar")
	defer def()

	deployment := Deployment(tmp)
	branches, err := deployment.ListBranches("proj")
	if err != nil {
		t.Fatal(err.Error())
	}

	expect := []string{"bar", "fancy", "foo", "master"}

	if !reflect.DeepEqual(branches, expect) {
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Obtained: %#v", branches)
	}
}
