/*

Package parser extract information from go test JSON output.

It looks at the incoming JSON logs and parse them to extract the
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


Normal output:

	$ go test file_test.go
	--- FAIL: Test (0.00s)
			a_test.go:6: a
			a_test.go:7: b
	FAIL
	FAIL    command-line-arguments  0.149s
	FAIL

JSON Output:
	$ go test -json file_test.go | jq '"\(.Action) | \(.Output)"'

	Step | Action | output
	---  | ---    | ---
	1    | run    |
	2    | output | === RUN   Test
	3    | output | a_test.go:6: a
	3    | output | a_test.go:7: b
	4    | output | --- FAIL: Test (0.00s)
	5    | fail   |
	X    | output | FAIL
	X    | output | FAIL\tcommand-line-arguments\t0.105s
	X    | fail   |

1. Initialise a new test
2. Extract and store the test name
3. Store all the tests logs
4. Store whether the test was a success
5. pass/fail: End of test
X. Ignored: Global test output

Output:

Test {
	Name: "Test",
	Package: "path/to/file_test.go",
	Lines: []Line{
		{ File: "a_test.go" Number: "6", Output: "a" },
		{ File: "a_test.go" Number: "7", Output: "b" },
	},
	Output: ""
}

See comments "Step X" in function Parser for steps details.

*/
package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	// lineNumberReExp holds the regexp to parse the Step 4. output, it will extract
	// the Filename, Line number and First (sometimes only) log line. Which gives a
	// total of 4 outputs when we call FindStringSubmatch.
	// E.g. '    file_test.go:42: Error Message\n'
	//     ['    file_test.go:42: Error Message\n', 'file_test.go', '42', 'Error Message']
	lineNumberReExp  = `^[[:space:]]+([[:word:]]+_test.go):([[:digit:]]+):[[:space:]](.*)\n$`
	expectedMatchLen = 4
	// We're going to parse the first 4 chars from the Output to distinguish
	// between step 2, 3 and 6.
	prefixLength = 4
)

// JSONLine holds the information for one log line coming from `go test -json`.
type JSONLine struct {
	Time    string `json:",omitempty"`
	Action  string `json:",omitempty"`
	Package string `json:",omitempty"`
	Test    string `json:",omitempty"`
	Output  string `json:",omitempty"`
}

type Line struct {
	Number string
	Output string
}

// Test holds the information for one test output given by `go test -json`.
type Test struct {
	Failed bool
	Added  bool

	Package string // Package of the test
	File    string
	Name    string // Test name as in `func TestName(...){`

	Lines  []Line // Line-based output
	Output string // Unparsed output
}

// Parse read each line from the scanner and create a list of tests that are failing.
// This function is long because all the logic has to be held for an unknown
// number of log lines, may refactor later.
// nolint: gocognit
func Parse(s *bufio.Scanner) ([]Test, error) {
	lineNumberRe := regexp.MustCompile(lineNumberReExp)

	tests := []Test{}
	test := Test{}

	for s.Scan() {
		// Parse the input line.
		line := JSONLine{}
		if err := json.Unmarshal(s.Bytes(), &line); err != nil {
			return nil, err
		}

		// Switch on the action
		switch line.Action {

		// Step 1.
		// Create a new test for the current package
		case "run":
			test = Test{
				Failed:  false,
				Added:   false,
				Package: line.Package,
			}

		// Step 5.
		// pass/fail reports the end of the test.
		// Add the test to the test list.
		// If the test is failed, the test is defined and the current log line is
		// not empty, it means that we're in the case 7 above, we can insert the
		// test in the list.
		case "pass":
			// test.Failed = false
			// tests = append(tests, test)

		// Step 5. same
		case "fail":
			// if testFailed && test.Line != "" && line.Test != "" {
			test.Failed = true
			if !test.Added {
				test.Added = true
				tests = append(tests, test)
			}
			// }

		// Other action need to be parsed based on the output.
		case "output":
			// Sometimes the output doesn't contain 5 chars, in which case there's no
			// point in parsing it.
			if len(line.Output) <= prefixLength {
				continue
			}

			switch line.Output[0:prefixLength] {

			// Step 2.
			// Save the test name
			case "=== ":
				test.Name = line.Test

			// Ignore because we're already storing this in step 5.
			case "--- ":
			case "FAIL":
				// pass

			// Parse the message manually
			default:
				matches := lineNumberRe.FindStringSubmatch(line.Output)
				fmt.Println(line.Output)
				fmt.Println(matches)
				// Step 4. Match a file/line number and message.
				if len(matches) == expectedMatchLen {
					// XXX what's this
					// if test.Failed && test.Line != "" {
					// tests = append(tests, test)
					// }

					test.File = matches[1]
					// Add filename, line, log line in test object
					test.Lines = append(test.Lines, Line{
						Number: matches[2],
						Output: matches[3],
					})
				} else {
					test.Output += "\n" + strings.TrimRight(line.Output, "\n")
				}
			}
		}

	}

	return tests, nil
}

// func ParseOld(s *bufio.Scanner) ([]Test, error) {
// lineNumberRe := regexp.MustCompile(lineNumberReExp)

// // Need to hold this states in between logs.
// testFailed := false

// tests := []Test{}
// test := Test{}

// for s.Scan() {
// // Parse the input line.
// line := JSONLine{}
// if err := json.Unmarshal(s.Bytes(), &line); err != nil {
// return nil, err
// }

// // Switch on action will lead to 1,
// switch line.Action {

// // Step 1.
// // Create a new test for the current package
// case "run":
// test = Test{
// Package: line.Package,
// }
// testFailed = false

// // Step 7.
// // If the test is failed, the test is defined and the current log line is
// // not empty, it means that we're in the case 7 above, we can insert the
// // test in the list.
// case "pass":
// case "fail":
// if testFailed && test.Line != "" && line.Test != "" {
// tests = append(tests, test)
// }
// case "output":
// // Sometimes the output doesn't contain 5 chars, in which case there's no
// // point in parsing it.
// if len(line.Output) <= prefixLength {
// continue
// }

// switch line.Output[0:prefixLength] {

// // Step 2.
// // Store the test function
// case "=== ":
// test.Name = line.Test

// // Step 4.
// // Store whether the test failed.
// case "--- ":
// if strings.HasPrefix(line.Output, "--- FAIL") {
// testFailed = true
// }

// // Step 6.
// // We don't have anything to do here, all the information are fetched
// // from other steps.
// case "FAIL":
// // pass
// default:
// matches := lineNumberRe.FindStringSubmatch(line.Output)
// // Step 4. Match a file/line number and message.
// if len(matches) == expectedMatchLen {
// // XXX what's this
// if testFailed && test.Line != "" {
// tests = append(tests, test)
// }

// // Get filename, line, first log line
// test = Test{
// // Propagate information from the last test.
// Name:    test.Name,
// Package: test.Package,
// // Add new parsed information.
// File: matches[1],
// Line: matches[2],
// // Output: matches[3],
// }
// } else {
// // Step 5. Add all other messages
// // test.Output += "\n" + strings.TrimRight(line.Output, "\n")
// }
// }
// }
// }

// return tests, nil
// }
