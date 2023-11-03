package source

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

const (
	readme = "README.md"
)

func TestDirectory(t *testing.T) {
	dir := tempdir(t)
	defer os.RemoveAll(dir)

	src, err := NewDirectory(&url.URL{Scheme: "dir", Path: dir})
	if err != nil {
		t.Error("could not setup dir source:", err)
	}

	// Create the readme file in dir.
	createfile(t, filepath.Join(dir, readme), []byte(content))

	compare(t, src, readme, content)
}

func TestDirectoryWatch(t *testing.T) {
	dir := tempdir(t)
	defer os.RemoveAll(dir)

	src, err := NewDirectory(&url.URL{Scheme: "dir", Path: dir})
	if err != nil {
		t.Error("could not setup dir source:", err)
	}

	w := src.(Watcher)
	defer w.Close()
	files, errors := w.Watch()

	// Create the readme file in dir while the watcher is running.
	createfile(t, filepath.Join(dir, readme), []byte(content))

	select {
	case f := <-files:
		if f != readme {
			t.Error("unexpected file", f)
		}
	case err := <-errors:
		t.Error("error occured while watching:", err)
	}
}
