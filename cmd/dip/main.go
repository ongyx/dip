package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/ongyx/dip"
)

var (
	address = flag.String("address", ":8080", "address to host the server at")
	version = flag.Bool("version", false, "show version of dip")
)

func main() {
	flag.Parse()

	if *version {
		ver := "(how did we get here?)"
		if bi, ok := debug.ReadBuildInfo(); ok {
			ver = bi.Main.Version
		}

		fmt.Printf("dip %s\n", ver)
		os.Exit(0)
	}

	source, err := dip.NewStdin()
	if err != nil {
		fmt.Println("failed to read stdin:", err)
		os.Exit(1)
	}

	server := setup(source)
	mux := &LogHandler{
		log:     logger,
		handler: server.Mux(nil),
	}

	addr := *address
	// address is only the port without a colon, so add it.
	if _, err := strconv.Atoi(addr); err == nil {
		addr = ":" + addr
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("failed to parse address %s: %s\n", addr, err)
		os.Exit(1)
	}

	// make the hostname localhost so it is a valid URI.
	if host == "" {
		host = "localhost"
	}

	fmt.Printf("listening at http://%s:%s\n", host, port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("server exited - %s\n", err)
	}
}
