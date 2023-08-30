package document

import (
	"io/fs"
	"strings"
)

const (
	// Root represents the root document path in a source.
	// Certain sources only serve a single Markdown file and thus may check for this value to serve that file.
	Root = "."
)

var (
	markdownExtensions = []string{".md", ".markdown"}
)

// Source is a filesystem serving Markdown files.
type Source interface {
	fs.FS

	// Watch watches the filesystem for changed Markdown files and sends their paths.
	// Any errors encountered when watching should be sent as well.
	Watch(files chan<- string, errors chan<- error)

	// Close performs cleanup on the source.
	Close() error
}

// IsMarkdownFile checks if the path ends with a Markdown file extension.
func IsMarkdownFile(path string) bool {
	for _, ext := range markdownExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
