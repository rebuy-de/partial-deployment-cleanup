package kv

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/hashicorp/consul/api"
)

var (
	namespace = flag.String(
		"consul-namespace",
		"nginx/partial_deployment/",
		"Root namespace for Consul KV")
	agent = flag.String(
		"consul-agent",
		"localhost:8500",
		"")
)

type Client struct {
	client *api.Client
}

func New() (*Client, error) {
	config := api.DefaultConfig()
	config.Address = *agent

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c := new(Client)
	c.client = client

	return c, nil
}

func (c *Client) GetProjects() ([]string, error) {
	kv := c.client.KV()
	pairs, _, err := kv.List(*namespace, nil)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	for _, pair := range pairs {
		key := strings.TrimPrefix(pair.Key, *namespace)
		key = strings.Split(key, "/")[0]
		keys = append(keys, key)
	}

	return keys, nil
}

func (c *Client) GetDeployments(project string) (Deployments, error) {
	ns := *namespace + project + "/deployments"

	kv := c.client.KV()
	pairs, _, err := kv.List(ns, nil)
	if err != nil {
		return nil, err
	}

	deployments := make(Deployments, 0)
	for _, pair := range pairs {
		branch := strings.TrimPrefix(pair.Key, ns)
		branch = strings.TrimLeft(branch, "/")
		branch = strings.Split(branch, "/")[0]

		var deployment Deployment
		err = json.Unmarshal(pair.Value, &deployment)
		if err != nil {
			return nil, err
		}

		deployment.Project = project
		deployment.Branch = branch

		deployments = append(deployments, &deployment)
	}

	return deployments, nil

}

func (c *Client) RemoveDeployment(d *Deployment) error {
	key := *namespace + d.Project + "/deployments/" + d.Branch
	kv := c.client.KV()
	_, err := kv.Delete(key, nil)
	return err
}
