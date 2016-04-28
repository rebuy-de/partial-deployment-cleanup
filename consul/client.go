package consul

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	namespace Key
	client    *api.Client
}

func New(agent string, namespace Key) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = agent

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	c := new(Client)
	c.namespace = namespace
	c.client = client

	return c, nil
}

func (c *Client) GetProjects() ([]string, error) {
	kv := c.client.KV()
	pairs, _, err := kv.List(string(c.namespace), nil)
	if err != nil {
		return nil, err
	}

	keyMap := make(map[string]bool)
	for _, pair := range pairs {
		key := Key(pair.Key).Base(c.namespace).Get(0)
		keyMap[key] = true
	}

	keys := make([]string, 0)
	for key, _ := range keyMap {
		keys = append(keys, key)
	}

	return keys, nil
}

func (c *Client) GetDeployments(project string) (Deployments, error) {
	ns := c.namespace.Add(project).Add("deployments")

	kv := c.client.KV()
	pairs, _, err := kv.List(string(ns), nil)
	if err != nil {
		return nil, err
	}

	deployments := make(Deployments, 0)
	for _, pair := range pairs {
		branch := Key(pair.Key).Base(ns).Get(0)

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
	key := c.namespace.Add(d.Project).Add("deployments").Add(d.Branch)
	kv := c.client.KV()
	_, err := kv.Delete(string(key), nil)
	return err
}

func (c *Client) GetDistribution(project string) (Distribution, error) {
	key := c.namespace.Add(project).Add("distribution")
	pair, _, err := c.client.KV().Get(key.Clean(), nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, fmt.Errorf("Didn't find a Distribution in namespace %s.", c.namespace)
	}

	var distribution Distribution
	err = json.Unmarshal(pair.Value, &distribution)
	if err != nil {
		return nil, err
	}

	return distribution, nil
}
