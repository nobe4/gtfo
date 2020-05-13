/*

Package parser looks at the incoming JSON logs and parse them to extract the
following information:
- Test File
- Test Name
- Test Line
- Test Output

So we can build a modular output with the information, for instance:

file:line:output


During the parsing, we'll have two returned tests for the following case:

func Test (t *testing.T) {
	t.Log("a")   <- first
	t.Fatal("b") <- second
}

Description of JSON logs:

  Step | Action    | Output
  1    | run       |
  2    | output    | === RUN
  3    | output    | --- PASS/FAIL
  4    | output    | file_test.go:42
  5    | output    | messages
  6    | output    | FAIL ...
  7    | pass/fail |

See comments "Step X" below for steps details.

*/
package parser

import (
	"bufio"
	"encoding/json"
	"regexp"
	"strings"
)

const (
	// lineNumberReExp holds the regexp to parse the Step 4. output, it will extract
	// the Filename, Line number and First (sometimes only) log line. Which gives a
	// total of 4 outputs when we call FindStringSubmatch.
	// E.g. '    file_test.go:42: Error Message\n'
	lineNumberReExp  = `^[[:space:]]+([[:word:]]+_test.go):([[:digit:]]+):[[:space:]](.*)\n$`
	expectedMatchLen = 4
	// We're going to parse the first 4 chars from the Output to distinguish
	// between wtep 2, 3 and 6.
	prefixLength = 4
)

// Line holds the information for one log line coming from `go test -json`.
type Line struct {
	Time    string `json:",omitempty"`
	Action  string `json:",omitempty"`
	Package string `json:",omitempty"`
	Test    string `json:",omitempty"`
	Output  string `json:",omitempty"`
}

// Test holds the information for one test output given by `go test -json`.
type Test struct {
	Package string
	File    string
	Name    string
	Line    string
	Output  string
}

// Parse read each line from the scanner and create a list of tests that are failing.
// This function is long because all the logic has to be held for an unknown
// number of log lines, may refactor later.
func Parse(s *bufio.Scanner) ([]Test, error) {
	lineNumberRe := regexp.MustCompile(lineNumberReExp)

	// Need to hold this states in between logs.
	testFailed := false

	tests := []Test{}
	test := Test{}

	for s.Scan() {
		// Parse the input line.
		line := Line{}
		if err := json.Unmarshal(s.Bytes(), &line); err != nil {
			return nil, err
		}

		// Switch on action will lead to 1,
		switch line.Action {
		case "run":
			// Step 1.
			// Create a new test for the current package
			test = Test{
				Package: line.Package,
			}
			testFailed = false
		case "pass":
		case "fail":
			// Step 7.
			// If the test is failed, the test is defined and the current log line is
			// not empty, it means that we're in the case 7 above, we can insert the
			// test in the list.
			if testFailed && test.Line != "" && line.Test != "" {
				tests = append(tests, test)
			}
		case "output":
			// Sometimes the output doesn't contain 5 chars, in which case there's no
			// point in parsing it.
			if len(line.Output) <= prefixLength {
				continue
			}

			switch line.Output[0:prefixLength] {
			case "=== ":
				// Step 2.
				// Store the test function
				test.Name = line.Test
			case "--- ":
				// Step 3.
				// Store if the test failed only.
				if strings.HasPrefix(line.Output, "--- FAIL") {
					testFailed = true
				}
			case "FAIL":
				// Step 6.
				// We don't have anything to do here, all the information are fetched
				// from other steps.
			default:
				matches := lineNumberRe.FindStringSubmatch(line.Output)
				if len(matches) == expectedMatchLen {
					// Step 4.
					if testFailed && test.Line != "" {
						tests = append(tests, test)
					}

					// Get filename, line, first log line
					test = Test{
						// Propagate information from the last test.
						Name:    test.Name,
						Package: test.Package,
						// Add new parsed information.
						File:   matches[1],
						Line:   matches[2],
						Output: matches[3],
					}
				} else {
					// Step 5.
					test.Output += "\n" + strings.TrimRight(line.Output, "\n")
				}
			}
		}
	}

	return tests, nil
}
