package main

import (
	"log"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/filesystem"
)

func CleanupFilesystem(agent string, namespace consul.Key, path string) error {
	log.Printf("Cleaning up file system with these parameters:\n"+
		"    agent:     %#v\n"+
		"    namespace: %#v\n"+
		"    path:      %#v",
		agent, namespace, path)

	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Handling these projects: %+v\n", projects)

	fs := filesystem.Deployment(path)

	for _, project := range projects {
		directories, err := fs.GetBranches(project)
		if err == filesystem.ProjectDirectoryNotFound {
			log.Printf("Skipping project %#v, because project directory doesn't exists.", project)
			continue
		} else if err != nil {
			return err
		}

		branches, err := client.GetBranches(project)
		if err != nil {
			return err
		}

		distribution, err := client.GetDistribution(project)
		if err != nil {
			return err
		}

		for _, branch := range directories {
			if branch == "master" {
				log.Printf("Keep branch %s/%s, because it is master.",
					project, branch)
				continue
			}

			if branches.Contains(branch) {
				log.Printf("Keep branch %s/%s, because it is still stored in Consul",
					project, branch)
				continue
			}

			if distribution.Contains(branch) {
				log.Printf("Keep branch %s/%s, because it is still listed in the distibution",
					project, branch)
				continue
			}

			log.Printf("Deleting branch %s/%s, because there is no deployment in Consul",
				project, branch)
			fs.Remove(project, branch)
		}
	}

	return nil
}
