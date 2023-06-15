package document

import (
	"io"
	"log"

	"github.com/ongyx/dip/internal/source"
	"github.com/yuin/goldmark"
)

// Library is an index of documents backed by a source.
type Library struct {
	*Index

	src source.Source
	md  goldmark.Markdown

	files  chan string
	errors chan error

	watching bool
}

// NewLibrary creates a new library with the given source and Markdown converter.
func NewLibrary(src source.Source, md goldmark.Markdown) *Library {
	l := &Library{
		Index:  NewIndex(),
		src:    src,
		md:     md,
		files:  make(chan string),
		errors: make(chan error),
	}

	return l
}

// Reload refreshes the content of a document given its file path.
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

	d := l.Index.Add(file)

	// use file name as document title
	if d.title == "" {
		fi, err := f.Stat()
		if err != nil {
			return err
		}

		d.title = fi.Name()
	}

	return d.Convert(text, l.md)
}

// Watch begins watching the library's source for changes.
// Errors are logged with the given logger.
func (l *Library) Watch(logger *log.Logger) {
	if l.watching {
		return
	}

	go func() {
		for file := range l.files {
			// only reload file if its already in the library
			if !l.Index.Has(file) {
				continue
			}

			if err := l.Reload(file); err != nil {
				l.errors <- err
			} else {
				logger.Printf("reloaded %s", file)
			}
		}
	}()

	go func() {
		for err := range l.errors {
			logger.Printf("error: library: %s\n", err)
		}
	}()

	go l.src.Watch(l.files, l.errors)

	l.watching = true
}

// Close closes the library.
func (l *Library) Close() error {
	close(l.files)
	close(l.errors)

	return l.src.Close()
}
