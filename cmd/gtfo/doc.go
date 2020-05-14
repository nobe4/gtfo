/*

GTFO parses and allow to format output from go test.

It runs by reading all the JSON from STDIN, parsing them and creating a
structured dataset.

Usage:

	go test -json ./... | gtfo [flags]

Flags:

	-format
		assigns the format to apply to the output; see internal/formatter for more
	- h
		display help message

*/
package main
