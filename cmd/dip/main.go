package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
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

	source, err := newSource(args.Path)
	if err != nil {
		fmt.Printf("failed to read from %s: %s\n", args.Path, err)
		os.Exit(1)
	}

	mux := setupHandler(source)

	addr := args.Address
	if isPort(addr) {
		addr = ":" + addr
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("failed to parse address %s: %s\n", addr, err)
		os.Exit(1)
	}

	// make the host localhost so that the address can open directly in a web browser.
	if host == "" {
		host = "localhost"
	}

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("listening at http://%s:%s\n", host, port)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("listener: %s\n", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// wait for ctrl+c to shutdown.
	<-interrupt

	fmt.Println(" shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Printf("error shutting down server: %s\n", err)
		os.Exit(2)
	}
}
