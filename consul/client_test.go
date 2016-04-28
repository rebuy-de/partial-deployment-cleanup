package consul

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/rebuy-de/partial-deployment-cleanup/testutil"
)

func TestGetProjects(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.Set("nginx/partial-deployment/green-web/foo", "bar")
	srv.Set("nginx/partial-deployment/green-web/blub", "blubber")
	srv.Set("nginx/partial-deployment/blue-web/foo", "bar")

	client, err := New(srv.Addr(), "nginx/partial-deployment/")
	if err != nil {
		t.Fatal(err.Error())
	}

	projects, err := client.GetProjects()
	if err != nil {
		t.Fatal(err.Error())
	}

	expect := []string{"blue-web", "green-web"}
	sort.Sort(sort.StringSlice(projects))

	if !reflect.DeepEqual(expect, []string(projects)) {
		t.Errorf("Obtained: %#v\n", projects)
		t.Errorf("Expected: %#v\n", expect)
	}
}

func TestGetDeployments(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.Set("nginx/partial-deployment/green-web/deployments/master",
		`{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}`)

	client, err := New(srv.Addr(), "nginx/partial-deployment/")
	if err != nil {
		t.Fatal(err.Error())
	}

	branches, err := client.GetBranches("green-web")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(branches) != 1 {
		t.Logf("%+v", branches)
		t.Fatalf("Expected 1 branch, but got %d", len(branches))
	}

	branch := *branches[0]

	expect := Branch{}
	expect.Project = "green-web"
	expect.Name = "master"
	expect.CreatedAt = time.Unix(1335205543, 0)
	expect.UpdatedAt = time.Unix(1335205543, 0)

	if branch.Project != "green-web" ||
		branch.Name != "master" ||
		branch.CreatedAt.Unix() != 1335205543 ||
		branch.UpdatedAt.Unix() != 1335205543 {
		t.Logf("Actual: %#v", branch)
		t.Logf("Expect: %#v", expect)
		t.Fatalf("Deployment doesn't equals expected value")
	}
}

func TestRemoveDeployments(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.Set("nginx/partial-deployment/green-web/deployments/master",
		`{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}`)
	srv.Set("nginx/partial-deployment/green-web/deployments/fancy",
		`{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}`)

	client, err := New(srv.Addr(), "nginx/partial-deployment/")
	if err != nil {
		t.Fatal(err.Error())
	}

	branch := Branch{}
	branch.Project = "green-web"
	branch.Name = "fancy"
	client.RemoveBranch(&branch)

	keys := srv.List("nginx/partial-deployment")
	if len(keys) != 1 {
		t.Logf("%#v", keys)
		t.Fatalf("Expected 1 key, but got %d", len(keys))
	}

	if keys[0] != "nginx/partial-deployment/green-web/deployments/master" {
		t.Fatalf("Deleted the wrong branch.")
	}
}

func TestGetDistribution(t *testing.T) {
	srv := testutil.CreateServer(t)
	defer srv.Stop()

	srv.Set("nginx/partial-deployment/blue-web/distribution",
		`{
			"1":"fancy",
			"2":"fancy",
			"3":"fancy",
			"4":"master",
			"5":"master"
		}`)

	client, err := New(srv.Addr(), "nginx/partial-deployment")
	if err != nil {
		t.Fatal(err.Error())
	}

	distribution, err := client.GetDistribution("blue-web")
	if err != nil {
		t.Fatal(err.Error())
	}

	expect := Distribution{
		"1": "fancy",
		"2": "fancy",
		"3": "fancy",
		"4": "master",
		"5": "master",
	}

	if !reflect.DeepEqual(expect, distribution) {
		t.Errorf("Result does not have the expected value.")
		t.Errorf("Expected: %#v", expect)
		t.Errorf("Obtained: %#v", distribution)
		t.FailNow()
	}

}
