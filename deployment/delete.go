package deployment

import (
	"flag"
	"log"
	"os"
	"path"
)

var (
	deploymentPath = flag.String("deployment-path", "/opt/www", "")
)

func Delete(project, branch string) {
	directory := getPath(project, branch)

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

func getPath(project, branch string) string {
	return path.Join(*deploymentPath, project, branch)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
