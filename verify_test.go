package main

import (
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestVerifyCorrect(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("green-web", "fancy")
	tmp.Branch("green-web", "master")

	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Branch{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/fancy",
		consul.Branch{
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

	err := Verify(srv.Addr(), consul.Key("nginx/partial-deployment"), tmp.Path())
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestVerifyDistribution(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("green-web", "fancy")
	tmp.Branch("green-web", "master")

	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Branch{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/fancy",
		consul.Branch{
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
			"5":"missing"
		}`)

	err := Verify(srv.Addr(), consul.Key("nginx/partial-deployment"), tmp.Path())
	if err != VerificationFailed {
		t.Fatalf("Expected verification to fail. %#v", err)
	}
}

func TestVerifyDeployment(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("green-web", "fancy")
	tmp.Branch("green-web", "master")

	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/master",
		consul.Branch{
			"green-web",
			"master",
			time.Now(),
			time.Now().Add(-4 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/missing",
		consul.Branch{
			"green-web",
			"fancy",
			time.Now(),
			time.Now().Add(-8 * Week),
		})
	srv.SetJSON(
		"nginx/partial-deployment/green-web/deployments/fancy",
		consul.Branch{
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

	err := Verify(srv.Addr(), consul.Key("nginx/partial-deployment"), tmp.Path())
	if err != VerificationFailed {
		t.Fatalf("Expected verification to fail. %#v", err)
	}
}
