package sse

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	ping = &Event{
		Type: "ping",
		Data: []byte("this is a test ping"),
	}
	pingText = `event: ping
data: this is a test ping

`
)

func splitStream(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	idx := bytes.Index(data, []byte("\n\n"))
	if idx >= 0 {
		// Extra 2 bytes for the double LF.
		return idx + 2, data[:idx], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func TestStream(t *testing.T) {
	t.Log("creating SSE stream")
	s := NewStream()

	ts := httptest.NewServer(H2CWrap(s))
	defer ts.Close()

	t.Log("preparing to connect via client")
	client := &http.Client{Transport: H2CTransport()}

	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Error("failed to prepare client request:", err)
	}

	t.Log("connecting")
	resp, err := client.Do(req)
	if err != nil {
		t.Error("could not request SSE connection:", err)
	}
	defer resp.Body.Close()
	t.Log("connected")

	// Send event only after the client has connected.
	t.Log("sending event")
	s.Send(ping)

	sc := bufio.NewScanner(resp.Body)
	sc.Split(splitStream)

	sc.Scan()
	// Put back the newlines as it was stripped by the split.
	got := sc.Text() + "\n\n"
	if got != pingText {
		t.Errorf("event does not match: expected %s, got %s", pingText, got)
	}
}
