package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nobe4/gtfo/internal/formatter"
	"github.com/nobe4/gtfo/internal/module"
	"github.com/nobe4/gtfo/internal/parser"
)

const (
	defaultFormat = `{{.Path}}:{{.Line}}: {{.Output}}\n`
	usageMessage  = "" +
		`Usage of %[1]s:

  go test -json ./... | %[1]s [flags]

Flags:
  -h bool
  	print this help
`
)

func main() {
	// Parsing the format
	formatString := flag.String(
		"format",
		defaultFormat,
		"Format to apply on the logs.\nSee internal/formatter for more info.",
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageMessage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Parse the log lines.
	s := bufio.NewScanner(os.Stdin)

	tests, err := parser.Parse(s)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get the current module, default to an empty string.
	module, err := module.Get()
	if err != nil {
		log.Printf("Couldn't figure out the current module: %v", err)
	}

	// Create the formatter, using the template and the module.
	format, err := formatter.Prepare(*formatString, module)
	if err != nil {
		log.Fatalf("Couldn't create the formatter: %v", err)
	}

	// Simple printing of the found logs.
	for _, test := range tests {
		fmt.Print(format(test))
	}

	// If there are any test, exit badly, so the editor can pick it up.
	if len(tests) > 0 {
		os.Exit(1)
	}
}
