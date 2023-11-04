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
	_ Watcher = &Directory{}

	markdownExtensions = []string{".md", ".markdown"}
)

// Directory is a source that serves a directory on the filesystem.
type Directory struct {
	fs      fs.FS
	path    string
	watcher *fsnotify.Watcher
}

// NewDirectory creates a new directory source.
func NewDirectory(u *url.URL) (Source, error) {
	path, err := filepath.Abs(u.Path)
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
	if isMarkdownFile(path) {
		return d.fs.Open(path)
	}

	return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrInvalid}
}

func (d *Directory) Watch() (<-chan string, <-chan error) {
	files := make(chan string)

	go func() {
		for event := range d.watcher.Events {
			if isMarkdownFile(event.Name) && event.Has(fsnotify.Write) {
				// SAFETY: event.Name is never outside of d.path
				rel, err := filepath.Rel(d.path, event.Name)
				if err != nil {
					panic("directory: not relative to path")
				}

				files <- rel
			}
		}

		close(files)
	}()

	return files, d.watcher.Errors
}

func (d *Directory) Close() error {
	return d.watcher.Close()
}

func isMarkdownFile(path string) bool {
	pathExt := filepath.Ext(path)

	for _, ext := range markdownExtensions {
		if pathExt == ext {
			return true
		}
	}

	return false
}
