package dip

import (
	"errors"
	"log"
)

var (
	// Root represents the root path in a source.
	Root = ""

	// ErrPathNotFound indicates the path was not found in a source.
	ErrPathNotFound = errors.New("path not found in source")
)

// Source represents a source of Markdown documents.
type Source interface {
	// Title returns the name of the document at the path.
	Title(path string) string

	// Read reads the Markdown document from the path as raw bytes in UTF8.
	// The path may be the root path (dip.Root).
	//
	// ErrPathNotFound should be returned as the error if the path was not found.
	Read(path string) ([]byte, error)

	// Reload sends a path on the queue when the document associated with the path should be reloaded
	// (i.e the document source has changed).
	Reload(queue chan<- string)

	// Log adds the logger to the source for debugging or logging errors.
	Log(logger *log.Logger)
}
