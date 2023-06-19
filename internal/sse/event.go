package sse

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

const (
	// EOF is the end of file value for serializing an event.
	EOF = '\n'
)

// Event represents a server-sent event.
// The zero value of an event can be used as-is.
//
// Note that any fields added must not contain double newlines (\n\n) as it represents EOF,
// and may truncate the rest of the event when read by a client.
//
// Reference: https://html.spec.whatwg.org/multipage/server-sent-events.html
type Event struct {
	buf bytes.Buffer
}

// Comment adds a comment to the event.
// This is ignored by the browser, so it can be used to send a heartbeat to keep the SSE connection alive.
func (e *Event) Comment(comment string) *Event {
	e.field("", comment)
	return e
}

// Type adds a type to the event.
// Event listeners can be dispatched client side to process the event.
func (e *Event) Type(eventType string) *Event {
	e.field("event", eventType)
	return e
}

// Data adds a payload to the event.
func (e *Event) Data(data []byte) *Event {
	e.field("data", data)
	return e
}

// ID adds a unique ID to the event.
func (e *Event) ID(id string) *Event {
	e.field("id", id)
	return e
}

// Retry adds a reconnection time to the event.
// This tells the client how long to wait before reconnecing to the server.
func (e *Event) Retry(reconnect time.Duration) *Event {
	e.field("retry", reconnect.Milliseconds())
	return e
}

// WriteTo copies the event into a writer.
func (e *Event) WriteTo(w io.Writer) (n int64, err error) {
	n, err = io.Copy(w, &e.buf)
	if err != nil {
		return
	}

	w.Write([]byte{EOF})

	return
}

// Reset resets the event, clearing all fields.
func (e *Event) Reset() {
	e.buf.Reset()
}

func (e *Event) field(name string, value any) {
	fmt.Fprintf(&e.buf, "%s: %s\n", name, value)
}
