package source

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	// Interface check.
	_ Watcher = &File{}
)

// File is a source that serves a single file on the filesystem.
type File struct {
	path    string
	watcher *fsnotify.Watcher
}

// NewFile creates a new file source.
func NewFile(u *url.URL) (Source, error) {
	path, err := filepath.Abs(u.Path)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Watch the parent directory of the file.
	// If we watch the file directly, any changes to the file cause the watch to be lost.
	// This is usually due to an temporary edited copy being moved over the existing one
	// (which is what Vim does).
	//
	// See https://pkg.go.dev/github.com/fsnotify/fsnotify#hdr-Watching_files.
	watcher.Add(filepath.Dir(path))

	return &File{
		path:    path,
		watcher: watcher,
	}, nil
}

func (f *File) Open(path string) (fs.File, error) {
	err := &fs.PathError{Op: "open", Path: path}

	if !fs.ValidPath(path) {
		err.Err = fs.ErrInvalid
		return nil, err
	}

	if path != Root {
		err.Err = fs.ErrNotExist
		return nil, err
	}

	return os.Open(f.path)
}

func (f *File) Watch() (<-chan string, <-chan error) {
	files := make(chan string)

	go func() {
		defer close(files)

		for event := range f.watcher.Events {
			// Reload the file only if it's being written to.
			if event.Name == f.path && event.Has(fsnotify.Write) {
				files <- Root
			}
		}
	}()

	return files, f.watcher.Errors
}

func (f *File) Close() error {
	return f.watcher.Close()
}
