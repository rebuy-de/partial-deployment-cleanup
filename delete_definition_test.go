package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/consul/testserver"
)

func TestRemoveDeployments(t *testing.T) {
	srv, def := testserver.Create(t)
	defer def()

	consul.Agent = &srv.HTTPAddr

	age := 4 * Week

	oldDeployment := consul.Deployment{
		"green-web",
		"master",
		time.Now(),
		time.Now().Add(-1 * age),
	}
	b, err := json.Marshal(oldDeployment)
	if err != nil {
		t.Fatal(err.Error())
	}
	srv.SetKV("nginx/partial-deployment/green-web/deployments/master", b)

	currentDeployment := consul.Deployment{
		"green-web",
		"fancy",
		time.Now(),
		time.Now(),
	}
	b, err = json.Marshal(currentDeployment)
	if err != nil {
		t.Fatal(err.Error())
	}
	srv.SetKV("nginx/partial-deployment/green-web/deployments/fancy", b)

	DeleteOldBranchDefinitions()

	keys := srv.ListKV("nginx/partial-deployment")
	if len(keys) != 1 {
		t.Logf("%#v", keys)
		t.Fatalf("Expected 1 key, but got %d", len(keys))
	}

	if keys[0] != "nginx/partial-deployment/green-web/deployments/fancy" {
		t.Fatalf("Deleted the wrong deployment.")
	}
}
