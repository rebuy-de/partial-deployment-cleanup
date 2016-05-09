package main

import (
	"log"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
)

func CleanupConsul(agent string, namespace consul.Key) error {
	log.Printf("Cleaning up Consul with these parameters:\n"+
		"    agent:     %#v\n"+
		"    namespace: %#v",
		agent, namespace)

	client, err := consul.New(agent, namespace)
	if err != nil {
		return err
	}

	projects, err := client.GetProjects()
	if err != nil {
		return err
	}

	log.Printf("Handling these projects: %+v\n", projects)

	for _, project := range projects {
		deployments, err := client.GetBranches(project)
		if err != nil {
			return err
		}

		distribution, err := client.GetDistribution(project)
		if err != nil {
			return err
		}

		for _, deployment := range deployments {
			age := time.Since(deployment.UpdatedAt)

			if deployment.Name == "master" {
				log.Printf("Keep branch %s/%s, because it is master.",
					project, deployment.Name)
				continue
			}

			if cleanupThreshold > age {
				log.Printf("Keep branch %s/%s, because it is only %s old.",
					project, deployment.Name, age.String())
				continue
			}

			if distribution.Contains(deployment.Name) {
				log.Printf("Keep branch %s/%s, because it is still listed in the distibution.",
					project, deployment.Name)
				continue
			}

			log.Printf("Deleting branch %s/%s, because it is %s old and is not listed in the distribution, anymore.",
				project, deployment.Name, age.String())

			err = client.RemoveBranch(deployment)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
