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

func (c *Client) GetProjects() (Projects, error) {
	kv := c.client.KV()
	pairs, _, err := kv.List(string(c.namespace), nil)
	if err != nil {
		return nil, err
	}

	projects := make(Projects, 0)
	for _, pair := range pairs {
		project := Key(pair.Key).Base(c.namespace).Get(0)
		if !projects.Contains(project) {
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func (c *Client) GetBranches(project string) (Branches, error) {
	ns := c.namespace.Add(project).Add("deployments")

	kv := c.client.KV()
	pairs, _, err := kv.List(string(ns), nil)
	if err != nil {
		return nil, err
	}

	branches := make(Branches, 0)
	for _, pair := range pairs {
		name := Key(pair.Key).Base(ns).Get(0)

		var branch Branch
		err = json.Unmarshal(pair.Value, &branch)
		if err != nil {
			return nil, err
		}

		branch.Project = project
		branch.Name = name

		branches = append(branches, &branch)
	}

	return branches, nil
}

func (c *Client) RemoveBranch(branch *Branch) error {
	key := c.namespace.Add(branch.Project).Add("deployments").Add(branch.Name)
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
