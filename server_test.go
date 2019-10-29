package ginserver

import "testing"

func TestServer(t *testing.T) {
	e := newServer()
	if e == nil {
		t.Errorf("Server should not be nil")
	}
}
