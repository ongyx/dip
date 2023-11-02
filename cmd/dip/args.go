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
		Address: "localhost:8080",
		Flags:   flag.CommandLine,
	}
)

func init() {
	args.Define()
}

type Args struct {
	Path    string
	Address string

	Version bool

	Flags *flag.FlagSet
}

func (a *Args) Define() {
	a.Flags.BoolVar(&args.Version, "v", false, "")
	a.Flags.BoolVar(&args.Version, "version", false, "")

	a.Flags.Usage = args.Usage
}

func (a *Args) Parse() {
	flag.Parse()

	if path := flag.Arg(0); path != "" {
		a.Path = path
	}

	if address := flag.Arg(1); address != "" {
		a.Address = address
	}
}

func (a *Args) Usage() {
	fmt.Fprintln(a.Flags.Output(), help)
}
