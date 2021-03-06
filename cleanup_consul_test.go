package main

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestCleanupConsul(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Branch{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/old",
		consul.Branch{
			"green-web",
			"old",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/ancient",
		consul.Branch{
			"green-web",
			"ancient",
			time.Now(),
			time.Now().Add(-8 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/fancy:",
		consul.Branch{
			"green-web",
			"fancy",
			time.Now(),
			time.Now(),
		})

	srv.Set("nginx/partial-deployment/green-web/distribution",
		`{
			"1":"fancy",
			"2":"fancy",
			"3":"fancy",
			"4":"old",
			"5":"old"
		}`)

	err := CleanupConsul(srv.Addr(), consul.Key("nginx/partial-deployment"))
	if err != nil {
		t.Fatal(err.Error())
	}

	expect := []string{
		"nginx/partial-deployment/green-web/deployments/fancy",
		"nginx/partial-deployment/green-web/deployments/master",
		"nginx/partial-deployment/green-web/deployments/old",
	}
	obtain := srv.List("nginx/partial-deployment/green-web/deployments")
	sort.Sort(sort.StringSlice(obtain))

	if reflect.DeepEqual(expect, obtain) {
		t.Errorf("Deleted the wrong deployment.")
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Obtained: %#v", obtain)
	}
}
