package main

import (
	"flag"
	"log"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/kv"
)

const (
	Day  time.Duration = 24 * time.Hour
	Week time.Duration = 7 * Day
)

var (
	cleanupThreshold = 2 * Week
)

func main() {
	flag.Parse()
}

func DeleteOldBranchDefinitions() error {
	client, err := kv.New()
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
			if cleanupThreshold < age {
				log.Printf("Deleting branch %s from %s, because it is to old.",
					deployment.Branch, project)
				log.Printf("%#v\n", deployment)

				err = client.RemoveDeployment(deployment)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func Redistribute() {
}

func DeleteOldBranchDeployments() {
}
