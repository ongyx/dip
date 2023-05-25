package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
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

	addr := args.Address
	if isPort(addr) {
		// since it's a standalone port, make it into a proper TCP address
		addr = ":" + addr
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("error: failed to parse address %s: %s\n", addr, err)
		os.Exit(1)
	}

	if host == "" {
		host = "127.0.0.1"
	}

	handler, err := newHandler(args.Path)
	if err != nil {
		fmt.Printf("error: failed to setup server: %s\n", err)
		os.Exit(2)
	}

	server := &http.Server{Addr: addr, Handler: handler}

	fmt.Printf("serving %s at http://%s:%s\n", args.Path, host, port)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("error: listener: %s\n", err)
			os.Exit(3)
		}
	}()

	// wait for ctrl+c to shutdown.
	wait()

	fmt.Println(" shutting down...")

	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Printf("error: shutdown: %s\n", err)
		os.Exit(3)
	}
}
