package sse

import (
	"bufio"
	"bytes"
	"net/http"
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

type client struct {
	resp    *http.Response
	scanner *bufio.Scanner
}

func newClient(hc *http.Client, u string) (*client, error) {
	r, err := hc.Get(u)
	if err != nil {
		return nil, err
	}

	sc := bufio.NewScanner(r.Body)
	sc.Split(splitStream)

	return &client{r, sc}, nil
}

func (c *client) nextEvent() string {
	if c.scanner.Scan() {
		// Put back the newlines as it was stripped by the split.
		return c.scanner.Text() + "\n\n"
	}

	return ""
}

func (c *client) close() {
	c.resp.Body.Close()
}

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
