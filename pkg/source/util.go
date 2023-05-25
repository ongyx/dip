package source

import "strings"

var (
	markdownExtensions = []string{".md", ".markdown"}
)

// IsMarkdownFile checks if the file path has a Markdown extension.
func IsMarkdownFile(path string) bool {
	for _, ext := range markdownExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
