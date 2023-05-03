package main

import (
	"flag"
	"fmt"
)

var (
	help = fmt.Sprintf(
		`dip: document instant preview for Markdown

        usage:
            dip [options] [<path>] [<address>]

        where:
            path is the file or directory to serve, or '-' for stdin. (default: '%s')
            address is the TCP address and/or port to listen for requests from. (default: '%s')

        options:`,
		args.Path,
		args.Address,
	)

	args = &Args{
		Path:    ".",
		Address: ":8080",
		Version: flag.Bool("version", false, "show version"),
	}
)

func init() {
	flag.Usage = args.Usage
}

type Args struct {
	Path    string
	Address string

	Version *bool
}

func (a *Args) Parse() {
	flag.Parse()

	path := flag.Arg(0)
	if path != "" {
		a.Path = path
	}

	address := flag.Arg(1)
	if address != "" {
		a.Address = address
	}
}

func (a *Args) Usage() {
	fmt.Fprintln(flag.CommandLine.Output(), help)
	flag.PrintDefaults()
}
