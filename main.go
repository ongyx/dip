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

	srv, err := createServer(args.Path)
	if err != nil {
		fmt.Printf("error: failed to setup server: %s\n", err)
		os.Exit(2)
	}

	hsrv := &http.Server{Addr: addr, Handler: wrapServer(srv)}
	go func() {
		if err := hsrv.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("error: listen: %s\n", err)
		}
	}()

	fmt.Printf("serving %s at http://%s:%s\n", args.Path, host, port)

	// wait for ctrl+c to shutdown.
	wait()

	fmt.Println("shutting down...")

	if err := srv.Close(); err != nil {
		fmt.Printf("error: close: %s\n", err)
	}

	if err := hsrv.Shutdown(context.Background()); err != nil {
		fmt.Printf("error: shutdown: %s\n", err)
	}
}
