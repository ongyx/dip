package source

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// File is a source that reads from a single file.
type File struct {
	base    string
	path    string
	watcher *fsnotify.Watcher
}

// NewFile creates a new file source.
func NewFile(path string) (Source, error) {
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
		base:    filepath.Base(path),
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

	if path != f.base {
		err.Err = fs.ErrNotExist
		return nil, err
	}

	return os.Open(path)
}

func (f *File) Root() string {
	return f.base
}

func (f *File) Watch(files chan<- string, errors chan<- error) {
	for {
		select {
		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}

			if event.Name == f.path && event.Has(fsnotify.Write) {
				files <- f.base
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}

			errors <- err
		}
	}
}

func (f *File) Close() error {
	return f.watcher.Close()
}
