package main

import (
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/ongyx/dip/pkg/document"
	"github.com/ongyx/dip/pkg/source"
	"github.com/ongyx/dip/pkg/web"
)

var (
	markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
)

func newHandler(path string) (http.Handler, error) {
	src, err := newSource(path)
	if err != nil {
		return nil, err
	}

	library := document.NewLibrary(src, markdown)
	server := web.NewServer(library, logger)
	mux := server.Mux(nil)

	return &LogHandler{logger: logger, handler: mux}, err
}

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

func isPort(addr string) bool {
	_, err := strconv.Atoi(addr)
	return err == nil
}

func wait() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
}
