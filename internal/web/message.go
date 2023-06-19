package web

// Message represents a document as an SSE payload.
type Message struct {
	// Content is the contents of the document's buffer.
	Content string `json:"content"`

	// Timestamp is when the document was last reloaded.
	Timestamp int64 `json:"timestamp"`
}
