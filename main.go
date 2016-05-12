package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/SocialCodeInc/go-gelf/gelf"
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

func wrapError(err error) cli.ExitCoder {
	if err != nil {
		msg := fmt.Sprintf("ERROR - (%T) %s", err, err.Error())
		return cli.NewExitError(msg, 3)
	}

	return nil
}

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
		cli.StringFlag{
			Name:  "path",
			Value: "/opt/www",
			Usage: "path for the deployment directory",
		},
		cli.StringFlag{
			Name:  "graylog-address",
			Usage: "address of the Graylog server",
		},
		cli.BoolFlag{
			Name:  "quiet",
			Usage: "reduce log output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "consul",
			Usage: "cleanup old deployments from Consul",
			Action: func(c *cli.Context) error {
				agent := c.GlobalString("agent")
				namespace := c.GlobalString("namespace")

				err := CleanupConsul(agent, consul.Key(namespace))
				return wrapError(err)
			},
		},
		{
			Name:  "filesystem",
			Usage: "cleanup old deployments from filesystem",
			Action: func(c *cli.Context) error {
				agent := c.GlobalString("agent")
				namespace := c.GlobalString("namespace")
				path := c.GlobalString("path")

				err := CleanupFilesystem(agent, consul.Key(namespace), path)
				return wrapError(err)
			},
		},
		{
			Name:  "verify",
			Usage: "verify Consul configuration against filesystem",
			Action: func(c *cli.Context) error {
				agent := c.GlobalString("agent")
				namespace := c.GlobalString("namespace")
				path := c.GlobalString("path")

				err := Verify(agent, consul.Key(namespace), path)
				if err == VerificationFailed {
					fmt.Printf("CRITICAL - %s\n", err.Error())
					os.Exit(2)
				} else if err != nil {
					fmt.Printf("ERROR - %s\n", err.Error())
					os.Exit(3)
				}

				fmt.Printf("OK - Existing directories match Consul state\n")

				return nil
			},
		},
	}

	app.Before = func(c *cli.Context) error {
		writers := make([]io.Writer, 0)

		if !c.GlobalBool("quiet") {
			writers = append(writers, os.Stdout)
		}

		if c.IsSet("graylog-address") {
			addr := c.GlobalString("graylog-address")
			gelfWriter, err := gelf.NewWriter(addr)
			if err != nil {
				return err
			}
			writers = append(writers, gelfWriter)
		}

		log.SetOutput(io.MultiWriter(writers...))

		return nil
	}

	app.Run(os.Args)
}
