package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/Bradbev/plywood/plywood"
)

const help = `
Plywood is a log combiner.  It accepts as input multiple files and produces a single
output stream on stdout.  Each line of input should start with a timestamp so that 
logs can be sorted.  The first line of a log file determines the timestamp format for 
the entire file.  If a timestamp is not understood, then that entire file will be at the
beginning of output, and each line will start with 'z'.
If the environment variable PLYWOOD is set, then that file name will be loaded to
enable custom runtime formatters.  
The -showDefaults flag will print the default formatters in the correct format to be loaded.
See https://golang.org/pkg/time/#Parse for details on the time layout format.
When defining custom formatters you must provide a sample line and the expected output.
Plywood will always format the output time in UTC.


Flags:
 -a   Show absolute timestamp.  Defaults to false
 -r   Hide relative timestamp.  Defaults to false
 -showDefaults  Print the default timestamp formats

Usage:
plywood one.log two.log > plywood.log
`

var absolute = flag.Bool("a", false, "Enable absolute timestamps in output")
var hideRelative = flag.Bool("r", false, "Hide relative timestamps in output")
var showDefaults = flag.Bool("showDefaults", false, "Print known timestamps")

func main() {
	flag.Parse()

	if *showDefaults {
		plywood.PrintDefaults()
		os.Exit(0)
	}

	ply := plywood.DefaultPlywood()

	if customFile := os.Getenv("PLYWOOD"); customFile != "" {
		data, err := ioutil.ReadFile(customFile)
		check(err)
		var formatters []plywood.CustomFormat
		err = json.Unmarshal(data, &formatters)
		check(err)
		ply = plywood.PlywoodFromFormatters(formatters)
	}

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

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
