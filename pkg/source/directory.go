package source

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Directory is a source that serves a directory on the filesystem.
type Directory struct {
	fs      fs.FS
	path    string
	watcher *fsnotify.Watcher
}

// NewDirectory creates a new directory source.
func NewDirectory(path string) (Source, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher.Add(path)

	return &Directory{
		fs:      os.DirFS(path),
		path:    path,
		watcher: watcher,
	}, nil
}

func (d *Directory) Open(path string) (fs.File, error) {
	return d.fs.Open(path)
}

func (d *Directory) Watch(files chan<- string, errors chan<- error) {
	for {
		select {
		case event, ok := <-d.watcher.Events:
			if !ok {
				return
			}

			if IsMarkdownFile(event.Name) && event.Has(fsnotify.Write) {
				// SAFETY: event.Name is never outside of d.path
				rel, err := filepath.Rel(d.path, event.Name)
				if err != nil {
					panic("directory: not relative to path")
				}

				files <- rel
			}
		case err, ok := <-d.watcher.Errors:
			if !ok {
				return
			}

			errors <- err
		}
	}
}

func (d *Directory) Close() error {
	return d.watcher.Close()
}
