package testserver

import (
	"bytes"
	"io"
	"testing"

	"github.com/hashicorp/consul/testutil"
)

func Create(t *testing.T) (*testutil.TestServer, func()) {
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

	if t.Skipped() {
		t.FailNow()
	}

	return srv, func() {
		srv.Stop()
	}
}
