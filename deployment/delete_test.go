package deployment

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func testHelperTempDir(t *testing.T, project string, branch string) (string, func()) {
	tmp, err := ioutil.TempDir("", "pdtest")
	if err != nil {
		t.Logf(err.Error())
		t.FailNow()
	}

	os.MkdirAll(path.Join(tmp, project, branch), 0755)

	return tmp, func() {
		os.RemoveAll(tmp)
	}
}

func testHelperDeploymentDirectory(t *testing.T, project string, branch string) (string, func()) {
	tmp, fn := testHelperTempDir(t, project, branch)

	for _, sub := range []string{"current", "release", "shared"} {
		os.Mkdir(path.Join(tmp, project, branch, sub), 0755)
	}

	return tmp, fn
}

func TestDeploymentDelete(t *testing.T) {
	tmp, fn := testHelperDeploymentDirectory(t, "proj", "fancy")
	defer fn()

	deploymentPath = &tmp
	Delete("proj", "fancy")

	dir := path.Join(tmp, "proj", "fancy")
	if isDirectory(dir) {
		t.Logf("`Delete` should have deleted '%s', but didn't.", dir)
		t.Fail()
	}
}

func TestDeploymentDeleteSkipMaster(t *testing.T) {
	tmp, fn := testHelperDeploymentDirectory(t, "proj", "master")
	defer fn()

	deploymentPath = &tmp
	Delete("proj", "master")

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
		deploymentPath = &tmp
		Delete("proj", "fancy")

		dir := path.Join(tmp, "proj", "fancy")
		if !isDirectory(dir) {
			t.Logf("`Delete` shouldn't have deleted '%s', but did.", dir)
			t.Fail()
		}
	}
}
