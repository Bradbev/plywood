package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Bradbev/plywood/plywood"
)

const help = `
Plywood is a log combiner.  It accepts as input multiple files and produces a single
output stream on stdout.  Each line of input should start with a timestamp so that 
logs can be sorted.  The first line of a log file determines the timestamp format for 
the entire file.  If a timestamp is not understood, then that entire file will be at the
beginning of output, and each line will start with 'z'.
Future releases of Plywood will allow custom time formats.

Flags:
 -a   Show absolute timestamp.  Defaults to false
 -r   Hide relative timestamp.  Defaults to false

Usage:
plywood one.log two.log > plywood.log
`

var absolute = flag.Bool("a", false, "Enable absolute timestamps in output")
var hideRelative = flag.Bool("r", false, "Hide relative timestamps in output")

func main() {
	flag.Parse()
	ply := plywood.DefaultPlywood()
	ply.IncludeAbsoluteTime = *absolute
	ply.IncludeRelativeTime = *hideRelative == false

	if len(flag.Args()) == 0 {
		usage()
	}

	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("Cannot add %v, %v\n", file, err)
			continue
		}
		defer f.Close()

		ply.AddReader(file, f)
	}
	io.Copy(os.Stdout, ply)
}

func usage() {
	fmt.Println(help)
	os.Exit(1)
}
