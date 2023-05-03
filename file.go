package dip

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// File is a source that reads from a single file.
type File struct {
	// Logger to log errors to.
	Log *log.Logger

	name, path string
	watcher    *fsnotify.Watcher
}

// NewFile creates a new file source.
func NewFile(path string) (*File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher.Add(filepath.Dir(path))

	return &File{
		name:    filepath.Base(path),
		path:    path,
		watcher: watcher,
	}, nil
}

func (f *File) Title(path string) string {
	return f.name
}

func (f *File) Read(path string) ([]byte, error) {
	if path != Root {
		return nil, ErrPathNotFound
	}

	return os.ReadFile(f.path)
}

func (f *File) Reload(queue chan<- string) {
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}

			if event.Name == f.path && event.Has(fsnotify.Write) {
				queue <- Root
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}

			f.Log.Printf("error: watcher failed: %s\n", err)
		}
	}
}
