package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	help = strings.ReplaceAll(fmt.Sprintf(
		`(Markdown) Document instant preview
        
        Usage:
          dip [options] [<path>] [<address>]
               
        Where:
          path is the file or directory to serve, or '-' for stdin. (default: '%s')
          address is the TCP address and/or port to listen for requests from. (default: '%s')
        
        Options:
          -v, -version  Print dip's version and exit.
          -h, -help     Print this message and exit.`,
		args.Path,
		args.Address,
	), "\n        ", "\n")

	args = &Args{
		Path:    ".",
		Address: ":8080",
	}
)

func init() {
	flag.BoolVar(&args.Version, "v", false, "")
	flag.BoolVar(&args.Version, "version", false, "")

	flag.Usage = args.Usage
}

type Args struct {
	Path    string
	Address string

	Version bool
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
}
