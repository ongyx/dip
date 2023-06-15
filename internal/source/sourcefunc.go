package source

var (
	sources = make(map[string]SourceFunc)
)

// SourceFunc represents a function that creates a source from the given path.
type SourceFunc func(path string) (Source, error)

// Register adds the source function by a scheme.
// If the scheme is already registered, this is a no-op.
func Register(scheme string, fn SourceFunc) {
	if _, ok := sources[scheme]; !ok {
		sources[scheme] = fn
	}
}

// Available returns a slice of all registered source schemes.
func Available() []string {
	s := make([]string, 0, len(sources))
	for k := range sources {
		s = append(s, k)
	}

	return s
}
