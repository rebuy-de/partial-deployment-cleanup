package main

import (
	"log"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
)

func DeleteOldBranchDefinitions(agent string, namespace consul.Key) error {
	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Cleaning up these projects: %#v\n", projects)

	for _, project := range projects {
		deployments, err := client.GetDeployments(project)
		if err != nil {
			return err
		}

		for _, deployment := range deployments {
			age := time.Since(deployment.UpdatedAt)
			log.Printf("Branch %s/%s is %s old.", project, deployment.Branch, age.String())
			if cleanupThreshold < age {
				log.Printf("Deleting branch %s/%s, because it is too old.",
					project, deployment.Branch)

				err = client.RemoveDeployment(deployment)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
