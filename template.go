package dip

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
)

var (
	//go:embed template
	templateFS embed.FS

	templateFunctions = template.FuncMap{
		"TrimSlash":  func(s string) string { return strings.Trim(s, "/") },
		"StatusText": http.StatusText,
	}
	tmpl = template.Must(template.New("template").Funcs(templateFunctions).ParseFS(templateFS, "template/*.html"))
)
