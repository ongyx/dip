package document

import (
	"testing"

	"github.com/yuin/goldmark"
)

var (
	converter = goldmark.New()
)

func TestDocumentConvert(t *testing.T) {
	d := NewDocument(converter)

	if err := d.Convert([]byte("# Hello World!")); err != nil {
		t.Error("failed to convert: ", err)
	}

	d.Borrow(func(buf []byte) error {
		got := string(buf)
		expected := "<h1>Hello World!</h1>\n"

		if got != expected {
			t.Errorf("content does not match: expected '%s', got '%s'", expected, got)
		}

		return nil
	})
}
