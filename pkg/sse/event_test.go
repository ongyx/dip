package sse

import (
	"bytes"
	"testing"
)

type testEvent struct {
	event     Event
	marshaled string
}

func TestEventMarshal(t *testing.T) {
	testEvents := []testEvent{
		{
			Event{Type: "greeting", Data: []byte("Hello World!")},
			`event: greeting
data: Hello World!

`,
		},
		{
			Event{Data: []byte("this is a message")},
			`data: this is a message

`,
		},
		{
			Event{Comment: "sse moment", Type: "moment", ID: "1", Data: []byte(`{"timestamp":0}`)},
			`: sse moment
event: moment
data: {"timestamp":0}
id: 1

`,
		},
	}

	var buf bytes.Buffer

	for _, te := range testEvents {
		buf.Reset()
		te.event.Marshal(&buf)

		expected := te.marshaled
		got := buf.String()

		if got != expected {
			t.Fatalf("got %s, wanted %s\n", got, expected)
		}
	}
}
