package document

import (
	"io"

	"github.com/yuin/goldmark"

	"github.com/ongyx/dip/internal/source"
)

// Library is an index of documents backed by a source.
type Library struct {
	*Index

	src source.Source
	md  goldmark.Markdown

	files    chan string
	errors   chan error
	watching bool
}

// NewLibrary creates a new library with the given source and Markdown converter.
func NewLibrary(src source.Source, md goldmark.Markdown) *Library {
	l := &Library{
		src:    src,
		md:     md,
		Index:  NewIndex(),
		files:  make(chan string),
		errors: make(chan error),
	}

	return l
}

// Root returns the path that should be opened instead of the root path.
func (l *Library) Root() string {
	return l.src.Root()
}

// Reload refreshes the content of a document at path and returns the document.
func (l *Library) Reload(file string) error {
	f, err := l.src.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	text, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	l.Index.Add(file)

	return l.Index.Borrow(file, func(d *Document) error {
		if d.Title == "" {
			fi, err := f.Stat()
			if err != nil {
				return err
			}

			// use file name as document title
			d.Title = fi.Name()
		}

		return d.Convert(text, l.md)
	})
}

// Watch watches the library's source for changes to files and reloads them.
// The paths of files reloaded and any errors are sent over the respective channels.
func (l *Library) Watch() (files <-chan string, errors <-chan error) {
	if !l.watching {
		changed := make(chan string)

		go func() {
			for file := range changed {
				// only reload file if its already in the library
				if !l.Index.Has(file) {
					continue
				}

				if err := l.Reload(file); err != nil {
					l.errors <- err
				} else {
					l.files <- file
				}
			}
		}()

		go l.src.Watch(changed, l.errors)

		l.watching = true
	}

	return l.files, l.errors
}

// Close closes the library.
func (l *Library) Close() error {
	close(l.files)
	close(l.errors)

	return l.src.Close()
}
