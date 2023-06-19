package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"github.com/ongyx/dip/internal/document"
	"github.com/ongyx/dip/internal/source"
	"github.com/ongyx/dip/internal/web"
)

var (
	markdown = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)
)

func wrapServer(h http.Handler) http.Handler {
	return web.H2Cify(&LogHandler{handler: h, logger: logger})
}

func createServer(path string) (*web.Server, error) {
	src, err := parseSource(path)
	if err != nil {
		return nil, err
	}

	lib := document.NewLibrary(src, markdown)
	return web.NewServer(lib, logger), nil
}

func parseSource(path string) (source.Source, error) {
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

	uri := scheme + "://" + path
	fmt.Printf("opening '%s'\n", uri)

	return source.New(uri)
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
