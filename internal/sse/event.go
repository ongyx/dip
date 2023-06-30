package sse

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	// EOL is the end-of-line character for fields and events.
	EOL = "\n"
)

// Event represents a server-sent event.
// The zero value of an event can be used as-is.
//
// Reference: https://html.spec.whatwg.org/multipage/server-sent-events.html
type Event struct {
	// Comment is a piece of text ignored by the client.
	// It may be used to send a heartbeat to keep the SSE connection alive.
	Comment string

	// Type is the kind of event.
	// Listeners for this specific type can be dispatched client side to process the event.
	Type string

	// Data is the event's text content.
	Data string

	// ID is a unique identifier for the event.
	ID string

	// Retry sets the client-side interval for reconnecting to the server.
	Retry time.Duration

	// If Raw is true, newlines in the event's data are not escaped.
	// Escaping prevents accidental truncation of text content.
	//
	// However, if it has already been encoded i.e. with JSON, Raw can be set to true.
	Raw bool
}

// Marshal returns the stream respresentation of the event.
func (e *Event) Marshal() []byte {
	var buf bytes.Buffer

	if e.Comment != "" {
		// escape newlines in comment anyway
		for _, line := range strings.Split(e.Comment, EOL) {
			marshal("", line, &buf)
		}
	}

	if e.Type != "" {
		marshal("event", e.Type, &buf)
	}

	if e.Data != "" {
		if e.Raw {
			marshal("data", e.Data, &buf)
		} else {
			for _, line := range strings.Split(e.Data, EOL) {
				marshal("data", line, &buf)
			}
		}
	}

	if e.ID != "" {
		marshal("id", e.ID, &buf)
	}

	if e.Retry != 0 {
		marshal("retry", strconv.FormatInt(e.Retry.Milliseconds(), 10), &buf)
	}

	buf.WriteString(EOL)

	return buf.Bytes()
}

func marshal(field, value string, w io.Writer) {
	fmt.Fprintf(w, "%s: %s" + EOL, field, value)
}
