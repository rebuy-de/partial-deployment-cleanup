package consul

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/consul/testutil"
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

func TestGetProjects(t *testing.T) {
	srv, def := testHelperConsul(t)
	defer def()

	Agent = &srv.HTTPAddr

	srv.SetKV("nginx/partial-deployment/green-web/foo", []byte("bar"))
	srv.SetKV("nginx/partial-deployment/green-web/blub", []byte("blubber"))
	srv.SetKV("nginx/partial-deployment/blue-web/foo", []byte("bar"))

	client, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	projects, err := client.GetProjects()
	if err != nil {
		t.Fatal(err.Error())
	}

	inArray := make(map[string]bool)
	for _, key := range projects {
		inArray[key] = true
	}

	t.Logf(`Want []string{"green-web", "blue-web"}`)

	if len(projects) != 2 {
		t.Fatalf("Got %#v\n", projects)
	}

	if _, ok := inArray["green-web"]; !ok {
		t.Fatalf("Got %#v\n", projects)
	}

	if _, ok := inArray["blue-web"]; !ok {
		t.Fatalf("Got %#v\n", projects)
	}
}

func TestGetDeployments(t *testing.T) {
	srv, def := testHelperConsul(t)
	defer def()

	Agent = &srv.HTTPAddr

	srv.SetKV("nginx/partial-deployment/green-web/deployments/master", []byte(`
		{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}
	`))

	client, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	deployments, err := client.GetDeployments("green-web")
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(deployments) != 1 {
		t.Logf("%+v", deployments)
		t.Fatalf("Expected 1 deployment, but got %d", len(deployments))
	}

	deployment := *deployments[0]

	expect := Deployment{}
	expect.Project = "green-web"
	expect.Branch = "master"
	expect.CreatedAt = time.Unix(1335205543, 0)
	expect.UpdatedAt = time.Unix(1335205543, 0)

	if deployment.Project != "green-web" ||
		deployment.Branch != "master" ||
		deployment.CreatedAt.Unix() != 1335205543 ||
		deployment.UpdatedAt.Unix() != 1335205543 {
		t.Logf("Actual: %#v", deployment)
		t.Logf("Expect: %#v", expect)
		t.Fatalf("Deployment doesn't equals expected value")
	}
}

func TestRemoveDeployments(t *testing.T) {
	srv, def := testHelperConsul(t)
	defer def()

	Agent = &srv.HTTPAddr

	srv.SetKV("nginx/partial-deployment/green-web/deployments/master", []byte(`
		{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}
	`))
	srv.SetKV("nginx/partial-deployment/green-web/deployments/fancy", []byte(`
		{
			"created_at": "2012-04-23T18:25:43Z",
			"updated_at": "2012-04-23T18:25:43Z"
		}
	`))

	client, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	deployment := Deployment{}
	deployment.Project = "green-web"
	deployment.Branch = "fancy"
	client.RemoveDeployment(&deployment)

	keys := srv.ListKV("nginx/partial-deployment")
	if len(keys) != 1 {
		t.Logf("%#v", keys)
		t.Fatalf("Expected 1 key, but got %d", len(keys))
	}

	if keys[0] != "nginx/partial-deployment/green-web/deployments/master" {
		t.Fatalf("Deleted the wrong deployment.")
	}
}
