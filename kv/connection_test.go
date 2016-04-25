package kv

import (
	"bytes"
	"io"
	"testing"

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

	agent = &srv.HTTPAddr

	srv.SetKV("nginx/partial_deployment/green-web/foo", []byte("bar"))
	srv.SetKV("nginx/partial_deployment/blue-web/foo", []byte("bar"))

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
