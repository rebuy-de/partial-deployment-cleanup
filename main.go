package main

import (
	"log"
	"os"
	"time"

	"github.com/codegangsta/cli"
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
	app := cli.NewApp()
	app.Name = "partial-deployment-cleanup"
	app.Usage = "Purges old branches and deployments"

	app.Commands = []cli.Command{
		{
			Name:    "consul",
			Aliases: []string{"a"},
			Usage:   "cleanup old deployments from Consul",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "namespace",
					Value: "nginx/partial-deployment",
					Usage: "Root namespace for Consul KV.",
				},
				cli.StringFlag{
					Name:  "agent",
					Value: "localhost:8500",
					Usage: "Host and port of the Consul agent. Should only be changed for development purposes.",
				},
			},
			Action: func(c *cli.Context) {
				agent := c.String("agent")
				namespace := c.String("namespace")

				err := CleanupConsul(agent, consul.Key(namespace))
				if err != nil {
					log.Print(err.Error())
					os.Exit(1)
				}
			},
		},
	}

	app.Run(os.Args)
}
