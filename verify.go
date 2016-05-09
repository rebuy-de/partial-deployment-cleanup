package main

import (
	"fmt"
	"log"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/filesystem"
)

var (
	VerificationFailed = fmt.Errorf("Verification failed!")
)

func Verify(agent string, namespace consul.Key, path string) error {
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

	failed := false

	for _, project := range projects {
		directories, err := fs.GetBranches(project)
		if err != nil {
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

		log.Printf("%#v has these branches on disk:     %+v", project, directories)
		log.Printf("%#v has these branches defined:     %+v", project, branches.Slice())
		log.Printf("%#v has these branches distributed: %+v", project, distribution.BranchSlice())

		for _, branch := range branches {
			if !directories.Contains(branch.Name) {
				log.Printf("%s/%s is defined on Consul, but don't exists on disk!", project, branch.Name)
				failed = true
			}
		}

		for _, branch := range distribution {
			if !directories.Contains(branch) {
				log.Printf("%s/%s is distributed in Consul, but don't exists on disk!", project, branch)
				failed = true
			}
		}
	}

	if failed {
		return VerificationFailed
	}

	return nil
}
