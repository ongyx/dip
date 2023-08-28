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
	EOL byte = '\n'
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

	// Data is the event's content.
	Data []byte

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
// If buf is not nil, the event is marshaled into the existing buffer.
func (e *Event) Marshal(buf *bytes.Buffer) []byte {
	if buf == nil {
		buf = new(bytes.Buffer)
	} else {
		buf.Reset()
	}

	if e.Comment != "" {
		// escape newlines in comment anyway
		for _, line := range strings.Split(e.Comment, string(EOL)) {
			marshal(buf, "", line)
		}
	}

	if e.Type != "" {
		marshal(buf, "event", e.Type)
	}

	if len(e.Data) > 0 {
		if e.Raw {
			marshal(buf, "data", e.Data)
		} else {
			for _, line := range bytes.Split(e.Data, []byte{EOL}) {
				marshal(buf, "data", line)
			}
		}
	}

	if e.ID != "" {
		marshal(buf, "id", e.ID)
	}

	if e.Retry != 0 {
		marshal(buf, "retry", strconv.FormatInt(e.Retry.Milliseconds(), 10))
	}

	buf.WriteByte(EOL)

	return buf.Bytes()
}

func marshal(w io.Writer, field string, value any) {
	fmt.Fprintf(w, "%s: %s%c", field, value, EOL)
}
