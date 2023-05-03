package dip

import (
	"log"

	"github.com/yuin/goldmark"
)

// Library parses documents from a source.
type Library struct {
	*Index

	// Logger for logging watcher event.
	Log *log.Logger

	// Path to static assets, i.e `/static/`.
	Static string

	// Markdown converter.
	Markdown goldmark.Markdown

	source Source
	queue  chan string
	event  chan Event
	init   bool
}

// NewLibrary creates a new library with the given source.
func NewLibrary(source Source) *Library {
	s := &Library{
		Index:  NewIndex(),
		source: source,
	}

	return s
}

// Title returns the title for the document at path.
func (l *Library) Title(path string) string {
	return l.source.Title(path)
}

// Reload reloads the document at path.
func (l *Library) Reload(path string) error {
	text, err := l.source.Read(path)
	if err != nil {
		return err
	}

	d := &Document{
		Title:  l.Title(path),
		Static: l.Static,
	}

	if err := d.Convert(text, l.Markdown); err != nil {
		return err
	}

	return l.Index.Cache(path, d)
}

// Watch watches the source for document changes and automatically reloads them.
// Any errors are sent over the event channel.
func (l *Library) Watch() <-chan Event {
	if !l.init {
		l.queue = make(chan string)
		l.event = make(chan Event)

		go func() {
			for path := range l.queue {
				err := l.Reload(path)

				l.event <- Event{Path: path, Error: err}
			}
		}()

		go l.source.Reload(l.queue)

		l.init = true
	}

	return l.event
}

// Close closes the library and stops watching the source for document changes.
func (l *Library) Close() {
	close(l.queue)
	close(l.event)
}

// Event represents the status of a document being reloaded.
type Event struct {
	Path  string
	Error error
}
