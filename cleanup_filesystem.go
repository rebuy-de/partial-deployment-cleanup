package main

import (
	"log"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/filesystem"
)

func CleanupFilesystem(agent string, namespace consul.Key, path string) error {
	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	deployment := filesystem.Deployment(path)

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Cleaning up these projects: %+v\n", projects)

	for _, project := range projects {
		deployed, err := deployment.ListBranches(project)
		if err != nil {
			return err
		}

		stored, err := client.GetDeployments(project)
		if err != nil {
			return err
		}

		storedSet := make(map[string]*consul.Deployment)
		for _, d := range stored {
			log.Printf("%#v", d.Branch)
			storedSet[d.Branch] = d
		}

		for _, branch := range deployed {
			if branch == "master" {
				log.Printf("Keep branch %s/%s, because it is master.",
					project, branch)
				continue
			}

			if _, ok := storedSet[branch]; ok {
				log.Printf("Keep branch %s/%s, because it is still stored in Consul",
					project, branch)
				continue
			}

			log.Printf("Deleting branch %s/%s, because there is no deployment in Consul",
				project, branch)
			deployment.Delete(project, branch)
		}
	}

	return nil
}
