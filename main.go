package main

import (
	"flag"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
)

const (
	Day  time.Duration = 24 * time.Hour
	Week time.Duration = 7 * Day
)

var (
	cleanupThreshold = 2 * Week
)

func main() {
	var (
		namespace = flag.String(
			"consul-namespace",
			"nginx/partial-deployment/",
			"Root namespace for Consul KV.")
		agent = flag.String(
			"consul-agent",
			"localhost:8500",
			"Host and port of the Consul agent. Should only be changed for development purposes.")
	)
	flag.Parse()

	DeleteOldBranchDefinitions(*agent, consul.Key(*namespace))
}

func DeleteOldBranchDeployments() {
}
