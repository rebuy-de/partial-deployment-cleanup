package main

import (
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestCleanupFilessystem(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Deployment{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/fancy:",
		consul.Deployment{
			"green-web",
			"fancy",
			time.Now(),
			time.Now().Add(-8 * Week),
		})
	srv.Set("nginx/partial-deployment/green-web/distribution",
		`{
			"1":"fancy",
			"2":"fancy",
			"3":"fancy",
			"4":"master",
			"5":"master"
		}`)

	err := CleanupConsul(srv.Addr(), consul.Key("nginx/partial-deployment"))
	if err != nil {
		t.Fatal(err.Error())
	}
}
