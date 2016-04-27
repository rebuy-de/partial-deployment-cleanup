package consul

import (
	"encoding/json"
	"strings"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	namespace string
	client    *api.Client
}

func New(agent string, namespace string) (*Client, error) {
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
	pairs, _, err := kv.List(c.namespace, nil)
	if err != nil {
		return nil, err
	}

	keyMap := make(map[string]bool)
	for _, pair := range pairs {
		key := strings.TrimPrefix(pair.Key, c.namespace)
		key = strings.Split(key, "/")[0]
		keyMap[key] = true
	}

	keys := make([]string, 0)
	for key, _ := range keyMap {
		keys = append(keys, key)
	}

	return keys, nil
}

func (c *Client) GetDeployments(project string) (Deployments, error) {
	ns := c.namespace + project + "/deployments"

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
	key := c.namespace + d.Project + "/deployments/" + d.Branch
	kv := c.client.KV()
	_, err := kv.Delete(key, nil)
	return err
}
