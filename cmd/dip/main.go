package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/ongyx/dip/pkg/document"
	"github.com/ongyx/dip/pkg/source"
)

func main() {
	args.Parse()

	if args.Version {
		ver := "(how did we get here?)"

		if bi, ok := debug.ReadBuildInfo(); ok {
			ver = bi.Main.Version
		}

		fmt.Printf("dip %s\n", ver)
		os.Exit(0)
	}

	if args.Path == "-" {
		fmt.Println("Type some Markdown below and press Ctrl-D (or Ctrl-Z + Enter on Windows).")
	}

	u, err := source.Parse(args.Path)
	if err != nil {
		fmt.Printf("error: could not parse path %s: %s\n", args.Path, err)
		os.Exit(1)
	}

	srv, err := document.NewServer(u, nil)
	if err != nil {
		fmt.Println("error: failed to setup server:", err)
		os.Exit(1)
	}
	srv.SetLogger(logger)

	// Pass address verbatim to http.Server.
	server := &http.Server{Addr: args.Address, Handler: wrap(srv)}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("error: listener:", err)
		}
	}()

	fmt.Printf("serving %s at http://%s\n", args.Path, args.Address)

	// wait for ctrl+c to shutdown.
	wait()

	fmt.Println("shutting down...")

	srv.Close()

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("error: shutdown:", err)
	}
}

func wait() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-interrupt
}

func wrap(h http.Handler) http.Handler {
	return h2c.NewHandler(&LogHandler{handler: h, logger: logger}, &http2.Server{})
}
