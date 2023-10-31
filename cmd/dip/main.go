package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
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

	src, err := source.Parse(args.Path)
	if err != nil {
		fmt.Printf("error: failed to get source for path %s - %s\n", args.Path, err)
		os.Exit(1)
	}

	library := document.NewLibrary(src, nil)
	server := document.NewServer(library, logger)

	httpServer := &http.Server{Addr: args.Address, Handler: wrap(server)}
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("error: listen: %s\n", err)
		}
	}()

	if strings.HasPrefix(args.Address, ":") {
		args.Address = "localhost" + args.Address
	}

	fmt.Printf("serving %s at http://%s\n", args.Path, args.Address)

	// wait for ctrl+c to shutdown.
	wait()

	fmt.Println("shutting down...")

	if err := httpServer.Shutdown(context.Background()); err != nil {
		fmt.Printf("error: shutdown: %s\n", err)
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
