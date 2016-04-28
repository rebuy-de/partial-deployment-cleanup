package main

import (
	"log"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
)

func CleanupConsul(agent string, namespace consul.Key) error {
	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Cleaning up these projects: %+v\n", projects)

	for _, project := range projects {
		deployments, err := client.GetDeployments(project)
		if err != nil {
			return err
		}

		distribution, err := client.GetDistribution(project)
		if err != nil {
			return err
		}

		for _, deployment := range deployments {
			age := time.Since(deployment.UpdatedAt)

			if deployment.Branch == "master" {
				log.Printf("Keep branch %s/%s, because it is master.",
					project, deployment.Branch)
				continue
			}

			if cleanupThreshold > age {
				log.Printf("Keep branch %s/%s, because it is only %s old.",
					project, deployment.Branch, age.String())
				continue
			}

			if distribution.Contains(deployment.Branch) {
				log.Printf("Keep branch %s/%s, because it is still listed in the distibution.",
					project, deployment.Branch)
				continue
			}

			log.Printf("Deleting branch %s/%s, because it is %s old and is not listed in the distribution, anymore.",
				project, deployment.Branch, age.String())

			err = client.RemoveDeployment(deployment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
