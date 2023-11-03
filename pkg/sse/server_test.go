package sse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewServer("")

	s.Add("test")
	if s.Get("test") == nil {
		t.Error("test stream does not exist")
	}

	// Removing non-existent streams should be fine.
	s.Remove("non-existent")

	s.Add("temp")
	s.Remove("temp")
	if s.Get("temp") != nil {
		t.Error("temp stream still exists")
	}
}

func TestServerSend(t *testing.T) {
	s := NewServer("")
	s.Add("test")

	ts := httptest.NewServer(H2CWrap(s))
	defer ts.Close()

	hc := &http.Client{Transport: H2CTransport()}

	c, err := newClient(hc, ts.URL+"?stream=test")
	if err != nil {
		t.Fatal("could not request SSE connection:", err)
	}
	defer c.close()

	s.Send("test", ping)

	got := c.nextEvent()
	if got != pingText {
		t.Errorf("event does not match: expected '%s', got '%s'", pingText, got)
	}
}

func TestServerSendDefault(t *testing.T) {
	s := NewServer("default")

	ts := httptest.NewServer(H2CWrap(s))
	defer ts.Close()

	hc := &http.Client{Transport: H2CTransport()}

	// Request the base URL without the 'stream' parameter.
	c, err := newClient(hc, ts.URL)
	if err != nil {
		t.Fatal("could not request SSE connection:", err)
	}
	defer c.close()

	s.Send("default", &Event{Type: "msg", Data: []byte("Welcome!")})

	expected := `event: msg
data: Welcome!

`
	got := c.nextEvent()
	if got != expected {
		t.Errorf("event does not match: expected '%s', got '%s'", expected, got)
	}
}
