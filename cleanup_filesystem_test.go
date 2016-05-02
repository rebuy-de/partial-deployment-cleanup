package main

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/consul"
	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestCleanupFilessystem(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	tmp := testutil.CreateDirectory(t)
	defer tmp.Clean()

	tmp.Branch("green-web", "fancy")
	tmp.Branch("green-web", "master")
	tmp.Branch("green-web", "old")
	tmp.Branch("green-web", "orphan")
	tmp.Branch("silo", "master")
	tmp.Branch("silo", "blubber")
	tmp.Branch("other", "master")
	tmp.Branch("other", "blubber")

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
			"5":"master",
			"6":"orphan"
		}`)
	srv.Set("nginx/partial-deployment/other/distribution",
		`{
			"1":"master",
			"2":"master",
			"3":"master",
			"4":"master",
			"5":"master"
		}`)

	err := CleanupFilesystem(srv.Addr(), consul.Key("nginx/partial-deployment"), tmp.Path())
	if err != nil {
		t.Fatal(err.Error())
	}

	globs, err := filepath.Glob(tmp.Path("*", "*"))
	if err != nil {
		t.Fatal(err.Error())
	}

	actual := make([]string, 0)
	for _, glob := range globs {
		short, err := filepath.Rel(tmp.Path(), glob)
		if err != nil {
			t.Fatal(err.Error())
		}
		actual = append(actual, short)
	}

	sort.Sort(sort.StringSlice(actual))

	expect := []string{
		"green-web/fancy",
		"green-web/master",
		"green-web/orphan",
		"other/master",
		"silo/blubber",
		"silo/master",
	}

	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("Got unexpected result.")
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Actual:   %#v", actual)
	}
}
