package source

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"testing"
)

const (
	nonexistent = "non-existent.md"
)

func TestFile(t *testing.T) {
	file := tempfile(t)
	name := file.Name()
	defer os.Remove(name)

	// The file must be closed to sync writes to the filesystem.
	file.WriteString(content)
	file.Close()

	src, err := NewFile(&url.URL{Scheme: "file", Path: name})
	if err != nil {
		t.Error("could not setup file source:", err)
	}

	// Try opening any file other than root.
	if _, err := src.Open(nonexistent); !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected non-root file to return fs.ErrNotExist, got", err)
	}

	compare(t, src, Root, content)
}

func TestFileWatch(t *testing.T) {
	file := tempfile(t)
	name := file.Name()
	defer os.Remove(name)

	src, err := NewFile(&url.URL{Scheme: "file", Path: name})
	if err != nil {
		t.Error("could not setup file source:", err)
	}

	w := src.(Watcher)
	defer w.Close()

	files, errors := w.Watch()

	// Write some content while the file watcher is running.
	// NOTE: The file needs to be closed for the file watcher to detect a write event!
	file.WriteString(content)
	file.Close()

	select {
	case f := <-files:
		if f != Root {
			t.Error("unexpected file", f)
		}
	case err := <-errors:
		t.Error("error occured while watching:", err)
	}
}
