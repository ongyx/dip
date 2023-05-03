package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/ongyx/dip"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	static = "__"
)

func newSource(path string) (dip.Source, error) {
	if path == "-" {
		return dip.NewStdin()
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return dip.NewDirectory(path)
	} else {
		return dip.NewFile(path)
	}
}

func setupHandler(source dip.Source) http.Handler {
	library := dip.NewLibrary(source)
	library.Static = static
	library.Markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	server := dip.NewServer(library, logger)

	return &LogHandler{
		log:     logger,
		handler: server.Mux(nil),
	}
}

func isPort(addr string) bool {
	_, err := strconv.Atoi(addr)
	return err == nil
}
