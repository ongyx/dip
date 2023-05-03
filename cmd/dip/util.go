package main

import (
	"net/http"
	"strconv"

	"github.com/ongyx/dip"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	static = "__"
)

func newSource(path string) (dip.Source, error) {
	if path == "-" {
		return dip.NewStdin()
	}

	return dip.NewFile(path)
}

func setupHandler(source dip.Source) http.Handler {
	library := dip.NewLibrary(source)
	library.Static = static
	library.Markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

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
