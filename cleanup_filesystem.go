package main

import (
	"log"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/filesystem"
)

func CleanupFilesystem(agent string, namespace consul.Key, path string) error {
	log.Printf("Cleaning up file system with these parameters:")
	log.Printf("    agent:     %#v", agent)
	log.Printf("    namespace: %#v", namespace)
	log.Printf("    path:      %#v", path)

	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Cleaning up these projects: %+v\n", projects)

	fs := filesystem.Deployment(path)

	for _, project := range projects {
		deployed, err := fs.GetBranches(project)
		if err != nil {
			return err
		}

		consulDeployments, err := client.GetBranches(project)
		if err != nil {
			return err
		}

		for _, branch := range deployed {
			if branch == "master" {
				log.Printf("Keep branch %s/%s, because it is master.",
					project, branch)
				continue
			}

			if consulDeployments.Contains(branch) {
				log.Printf("Keep branch %s/%s, because it is still stored in Consul",
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
