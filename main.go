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

func main() {
	flag.Parse()
	DeleteOldBranchDefinitions()
}

func Redistribute() {
}

func DeleteOldBranchDeployments() {
}
