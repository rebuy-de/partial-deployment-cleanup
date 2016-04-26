package main

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/consul/testutil"
	"github.com/rebuy-de/partial-deployment-cleanup/kv"
)

func testHelperConsul(t *testing.T) (*testutil.TestServer, func()) {
	var stdout io.Writer

	if !testing.Verbose() {
		stdout = &bytes.Buffer{}
	} else {
		stdout = nil
	}

	srv := testutil.NewTestServerConfig(t, func(c *testutil.TestServerConfig) {
		c.Server = true
		c.Bootstrap = true
		c.Stdout = stdout
	})

	return srv, func() {
		srv.Stop()
	}
}

func TestRemoveDeployments(t *testing.T) {
	srv, def := testHelperConsul(t)
	defer def()

	kv.Agent = &srv.HTTPAddr

	age := 4 * Week

	oldDeployment := kv.Deployment{
		"green-web",
		"master",
		time.Now(),
		time.Now().Add(-1 * age),
	}
	b, err := json.Marshal(oldDeployment)
	if err != nil {
		t.Fatal(err.Error())
	}
	srv.SetKV("nginx/partial_deployment/green-web/deployments/master", b)

	currentDeployment := kv.Deployment{
		"green-web",
		"fancy",
		time.Now(),
		time.Now(),
	}
	b, err = json.Marshal(currentDeployment)
	if err != nil {
		t.Fatal(err.Error())
	}
	srv.SetKV("nginx/partial_deployment/green-web/deployments/fancy", b)

	DeleteOldBranchDefinitions()

	keys := srv.ListKV("nginx/partial_deployment")
	if len(keys) != 1 {
		t.Logf("%#v", keys)
		t.Fatalf("Expected 1 key, but got %d", len(keys))
	}

	if keys[0] != "nginx/partial_deployment/green-web/deployments/fancy" {
		t.Fatalf("Deleted the wrong deployment.")
	}
}
