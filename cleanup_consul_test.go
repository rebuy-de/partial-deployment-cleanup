package main

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/consul/testserver"
)

func TestCleanupConsul(t *testing.T) {
	srv, def := testserver.Create(t)
	defer def()

	testserver.SetJSON(
		t, srv,
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Deployment{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	testserver.SetJSON(
		t, srv,
		"nginx/partial-deployment/green-web/deployments/old",
		consul.Deployment{
			"green-web",
			"old",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	testserver.SetJSON(
		t, srv,
		"nginx/partial-deployment/green-web/deployments/ancient",
		consul.Deployment{
			"green-web",
			"ancient",
			time.Now(),
			time.Now().Add(-8 * Week),
		})
	testserver.SetJSON(
		t, srv,
		"nginx/partial-deployment/green-web/deployments/fancy:",
		consul.Deployment{
			"green-web",
			"fancy",
			time.Now(),
			time.Now(),
		})

	srv.SetKV("nginx/partial-deployment/green-web/distribution", []byte(`
		{
			"1":"fancy",
			"2":"fancy",
			"3":"fancy",
			"4":"old",
			"5":"old"
		}
	`))

	err := CleanupConsul(srv.HTTPAddr, consul.Key("nginx/partial-deployment"))

	if err != nil {
		t.Fatal(err.Error())
	}

	expect := []string{
		"nginx/partial-deployment/green-web/deployments/fancy",
		"nginx/partial-deployment/green-web/deployments/master",
		"nginx/partial-deployment/green-web/deployments/old",
	}
	obtain := srv.ListKV("nginx/partial-deployment/green-web/deployments")
	sort.Sort(sort.StringSlice(obtain))

	if reflect.DeepEqual(expect, obtain) {
		t.Errorf("Deleted the wrong deployment.")
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Obtained: %#v", obtain)
	}
}
