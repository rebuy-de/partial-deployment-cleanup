package kv

import (
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
