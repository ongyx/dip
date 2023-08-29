package document

import "io/fs"

// Source is a filesystem serving Markdown files.
type Source interface {
	fs.FS

	// Watch watches the filesystem for changed Markdown files and sends their paths.
	// Any errors encountered when watching should be sent as well.
	Watch(files chan<- string, errs chan<- error)
}
