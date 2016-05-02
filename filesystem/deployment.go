package filesystem

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Deployment string

func (d *Deployment) Remove(project, branch string) {
	directory := path.Join(string(*d), project, branch)

	if branch == "master" {
		log.Printf(`Aborting deletion of directory '%s', `+
			`because it is the master branch.`, directory)
		return
	}

	for _, sub := range []string{"current", "release", "shared"} {
		if !isDirectory(path.Join(directory, sub)) {
			log.Printf(`Aborting deletion of directory '%s', `+
				`because it doesn't look like a deployment directory. `+
				`The subdirectory '%s' is missing.`, directory, sub)
			return
		}
	}

	os.RemoveAll(directory)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func (d *Deployment) GetBranches(project string) (Branches, error) {
	directory := path.Join(string(*d), project)

	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	branches := make([]string, 0)
	for _, fi := range fileInfos {
		branches = append(branches, fi.Name())
	}

	return branches, nil
}
