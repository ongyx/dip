package sse

import (
	"bytes"
	"testing"
)

type testMessage struct {
	message   Message
	marshaled string
}

func TestMessageMarshal(t *testing.T) {
	testMessages := []testMessage{
		{
			Message{"greeting", []byte("Hello World!")},
			`event: greeting
data: Hello World!

`,
		},
		{
			Message{"", []byte("this is a message")},
			`data: this is a message

`,
		},
	}

	var buf bytes.Buffer

	for _, tmsg := range testMessages {
		tmsg.message.Marshal(&buf)

		expected := tmsg.marshaled
		got := buf.String()

		if got != expected {
			t.Fatalf("got %s, wanted %s\n", got, expected)
		}
	}
}
