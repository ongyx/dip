package source

// Watcher is a source that supports file watching.
type Watcher interface {
	Source

	// Watch starts watching the source for Markdown file changes,
	// sending them over the files channel.
	// Any error encountered should be sent over the errors channel.
	//
	// When the watcher is closed, both returned channels should be closed.
	Watch() (files <-chan string, errors <-chan error)

	// Close closes the watcher.
	Close() error
}
