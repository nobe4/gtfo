package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nobe4/gtfo/internal/module"
	"github.com/nobe4/gtfo/internal/parser"
)

func main() {
	// Get the current module.
	module, err := module.Get()
	if err != nil {
		log.Fatalf("Couldn't figure out the current module: %v", err)
	}

	// Parse the log lines.
	s := bufio.NewScanner(os.Stdin)

	tests, err := parser.Parse(s)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Simple printing of the found logs.
	for _, test := range tests {
		// Remove the module from the package
		modifiedP := strings.TrimPrefix(test.Package, module+"/")

		fmt.Printf("%s:%s: %s\n", filepath.Join(modifiedP, test.File), test.Line, test.Output)
	}

	// If there are any test, exit badly, so the editor can pick it up.
	if len(tests) > 0 {
		os.Exit(1)
	}
}
