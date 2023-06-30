package sse

import (
	"testing"
)

type testEvent struct {
	event     Event
	marshaled string
}

func TestEventMarshal(t *testing.T) {
	testEvents := []testEvent{
		{
			Event{Type: "greeting", Data: "Hello World!"},
			`event: greeting
data: Hello World!

`,
		},
		{
			Event{Data: "this is a message"},
			`data: this is a message

`,
		},
		{
			Event{Comment: "sse moment", Type: "moment", ID: "1", Data: `{"timestamp":0}`},
			`: sse moment
event: moment
data: {"timestamp":0}
id: 1

`,
		},
	}

	for _, te := range testEvents {
		expected := te.marshaled
		got := string(te.event.Marshal())

		if got != expected {
			t.Fatalf("got %s, wanted %s\n", got, expected)
		}
	}
}
