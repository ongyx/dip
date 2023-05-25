package sse

import (
	"bytes"
	"fmt"
)

// Message represents a server-sent event.
type Message struct {
	// The event type. If empty, the event field is omitted.
	Event string

	// The data to send.
	Data []byte
}

// Marshal serializes the message into a buffer.
// Note that the buffer is reset beforehand.
func (m Message) Marshal(buf *bytes.Buffer) {
	buf.Reset()

	if m.Event != "" {
		buf.WriteString(fmt.Sprintf("event: %s\n", m.Event))
	}

	buf.WriteString(fmt.Sprintf("data: %s\n\n", m.Data))
}
