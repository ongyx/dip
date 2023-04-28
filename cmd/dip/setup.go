package main

import (
	"log"
	"os"

	"github.com/ongyx/dip"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var (
	static = "__"

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmsgprefix)
)

func setup(source dip.Source) *dip.Server {
	library := dip.NewLibrary(source)
	library.Static = static
	library.Markdown = goldmark.New(goldmark.WithExtensions(extension.GFM))

	return dip.NewServer(library, logger)
}
