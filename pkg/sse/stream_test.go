package sse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStream(t *testing.T) {
	t.Log("creating SSE stream")
	s := NewStream()

	ts := httptest.NewServer(H2CWrap(s))
	defer ts.Close()

	t.Log("preparing to connect via client")
	hc := &http.Client{Transport: H2CTransport()}

	t.Log("connecting")

	c, err := newClient(hc, ts.URL)
	if err != nil {
		t.Fatal("could not request SSE connection:", err)
	}
	defer c.close()

	t.Log("connected")

	// Send event only after the client has connected.
	t.Log("sending event")
	s.Send(ping)

	got := c.nextEvent()
	if got != pingText {
		t.Errorf("event does not match: expected '%s', got '%s'", pingText, got)
	}
}
