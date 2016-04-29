package testutil

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type TestDirectory string

func CreateDirectory(t *testing.T) TestDirectory {
	tmp, err := ioutil.TempDir("", "pdtest")
	if err != nil {
		t.Logf(err.Error())
		t.FailNow()
	}

	return TestDirectory(tmp)
}

func (tmp *TestDirectory) Path(subs ...string) string {
	args := []string{string(*tmp)}
	args = append(args, subs...)
	return path.Join(args...)
}

func (tmp *TestDirectory) Clean() {
	os.RemoveAll(tmp.Path())
}

func (tmp *TestDirectory) Branch(project, branch string) {
	os.MkdirAll(tmp.Path(project, branch), 0755)
	for _, sub := range []string{"current", "release", "shared"} {
		os.Mkdir(path.Join(tmp.Path(), project, branch, sub), 0755)
	}
}
