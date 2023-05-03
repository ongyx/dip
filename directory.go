package dip

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Directory is a source that reads from a directory.
type Directory struct {
	path    string
	watcher *fsnotify.Watcher

	log *log.Logger
}

// NewDirectory creates a new directory source.
// path must be absolute.
func NewDirectory(path string) (*Directory, error) {
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
		path:    path,
		watcher: watcher,
	}, nil
}

func (d *Directory) Title(path string) string {
	return resolve(path)
}

func (d *Directory) Read(path string) ([]byte, error) {
	abs := filepath.Join(d.path, resolve(path))

	return os.ReadFile(abs)
}

func (d *Directory) Reload(queue chan<- string) {
	for {
		select {
		case event, ok := <-d.watcher.Events:
			if !ok {
				return
			}

			if isMarkdown(event.Name) && event.Has(fsnotify.Write) {
				// SAFETY: event.Name is never outside of d.path
				rel, _ := filepath.Rel(d.path, event.Name)
				queue <- rel
			}
		case err, ok := <-d.watcher.Errors:
			if !ok {
				return
			}

			d.log.Printf("error: watcher: %s\n", err)
		}
	}
}

func (d *Directory) Log(logger *log.Logger) {
	d.log = logger
}

func resolve(path string) string {
	// special case: the root always redirects to README.md
	if path == Root {
		return "README.md"
	}

	return path
}
