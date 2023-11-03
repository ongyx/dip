package source

import (
	"io"
	"os"
	"testing"
)

func TestStdin(t *testing.T) {
	// Restore standard input later.
	stdin := os.Stdin
	defer func() {
		os.Stdin = stdin
	}()

	// Mock standard input using a pipe.
	r, w, err := os.Pipe()
	if err != nil {
		t.Error("failed to setup mock stdin:", err)
	}

	os.Stdin = r

	io.WriteString(w, content)
	w.Close()

	// Stdin never errors out.
	src, _ := NewStdin(nil)

	compare(t, src, Root, content)
}
