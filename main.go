package main

import (
	"fmt"
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
	version          = "unknown"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
	}

	app := cli.NewApp()
	app.Name = "partial-deployment-cleanup"
	app.Usage = "purges old branches and deployments"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "namespace",
			Value: "nginx/partial-deployment",
			Usage: "root namespace for Consul KV",
		},
		cli.StringFlag{
			Name:  "agent",
			Value: "localhost:8500",
			Usage: "host and port of the Consul agent. Should only be changed for development purposes",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "consul",
			Usage: "cleanup old deployments from Consul",
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
		{
			Name:  "filesystem",
			Usage: "cleanup old deployments from filesystem",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path",
					Usage: "path for the deployment directory",
				},
			},
			Action: func(c *cli.Context) {
				agent := c.String("agent")
				namespace := c.String("namespace")
				path := c.String("path")

				err := CleanupFilesystem(agent, consul.Key(namespace), path)
				if err != nil {
					log.Print(err.Error())
					os.Exit(1)
				}
			},
		},
	}

	app.Run(os.Args)
}
