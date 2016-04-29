package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	consulTestUtil "github.com/hashicorp/consul/testutil"
)

type TestServer struct {
	t   *testing.T
	srv *consulTestUtil.TestServer
}

func CreateServer(t *testing.T) *TestServer {
	var stdout io.Writer

	if !testing.Verbose() {
		stdout = &bytes.Buffer{}
	} else {
		stdout = nil
	}

	srv := consulTestUtil.NewTestServerConfig(t, func(c *consulTestUtil.TestServerConfig) {
		c.Server = true
		c.Bootstrap = true
		c.Stdout = stdout
	})

	if t.Skipped() {
		t.FailNow()
	}

	return &TestServer{t, srv}
}

func (s *TestServer) Addr() string {
	return s.srv.HTTPAddr
}

func (s *TestServer) Stop() {
	s.srv.Stop()
}

func (s *TestServer) Set(key string, value string) {
	s.srv.SetKV(key, []byte(value))
}

func (s *TestServer) List(prefix string) []string {
	return s.srv.ListKV(prefix)
}

func (s *TestServer) SetJSON(key string, value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		s.t.Fatal(err.Error())
	}
	s.srv.SetKV(key, data)
}
