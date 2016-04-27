package main

import (
	"flag"
	"time"
)

const (
	Day  time.Duration = 24 * time.Hour
	Week time.Duration = 7 * Day
)

var (
	cleanupThreshold = 2 * Week
)

var (
	namespace = flag.String(
		"consul-namespace",
		"nginx/partial-deployment/",
		"Root namespace for Consul KV")
	agent = flag.String(
		"consul-agent",
		"localhost:8500",
		"")
)

func main() {
	flag.Parse()
	DeleteOldBranchDefinitions(*agent, *namespace)
}

func Redistribute() {
}

func DeleteOldBranchDeployments() {
}
