package filesystem

import (
	"log"
	"os"
	"path"
)

type Deployment string

func (d *Deployment) Delete(project, branch string) {
	directory := d.getPath(project, branch)

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

func (d *Deployment) getPath(project, branch string) string {
	return path.Join(string(*d), project, branch)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
