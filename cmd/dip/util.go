package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/ongyx/dip/pkg/document"
	"github.com/ongyx/dip/pkg/source"
	"github.com/ongyx/dip/pkg/web"
)

var (
	static = "__"
)

func newSource(path string) (source.Source, error) {
	var scheme string

	if path == "-" {
		scheme = "stdin"
	} else {
		info, err := os.Stat(path)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			scheme = "dir"
		} else {
			scheme = "file"
		}
	}

	return source.New(scheme + "://" + path)
}

func setupHandler(src source.Source) http.Handler {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
	library := document.NewLibrary(src, md)

	server := web.NewServer(library, logger)

	return &LogHandler{
		log:     logger,
		handler: server.Mux(nil),
	}
}

func isPort(addr string) bool {
	_, err := strconv.Atoi(addr)
	return err == nil
}
