package source

// Watcher is a source that supports file watching.
type Watcher interface {
	Source

	// Watch watches the filesystem for changed Markdown files and sends their paths.
	// Any errors encountered when watching should be sent as well.
	Watch(files chan<- string, errors chan<- error)

	// Close performs cleanup on the source.
	Close() error
}
