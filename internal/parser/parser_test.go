package parser

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// Package used by `go test` when calling a single file
	pkg = "command-line-arguments"
)

// createScanner create a bufioScanner from the provided test.
// Returns the scanner and the filename for later test.
func createScanner(t *testing.T, s string) (*bufio.Scanner, string) {
	// Create the test content
	content := "package a"

	if len(s) > 0 {
		content += "\n" + `import "testing"`
	}

	content += "\n" + s

	tmpfile, err := ioutil.TempFile("", "*_test.go")
	assert.NoError(t, err)

	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	_, err = tmpfile.Write([]byte(content))
	assert.NoError(t, err)

	// Get the file name (not the path) for parsing check
	stat, err := tmpfile.Stat()
	assert.NoError(t, err)

	filename := stat.Name()

	// Go test exit with 1 when the tests are not passing, so we can safely ignore this error.
	// nolint:  gosec
	out, _ := exec.Command("go", "test", "-json", "-cover", "-race", tmpfile.Name()).Output()

	reader := strings.NewReader(string(out))

	return bufio.NewScanner(reader), filename
}

func TestCases(t *testing.T) {
	// In the tests, the line will be computed from the top of the file.
	// So the line 1 of the first test should be 5.
	testCases := []struct {
		description string
		testFunc    string
		expected    []Test
	}{
		{
			description: "No Test",
			testFunc:    ``,
			expected:    []Test{},
		},

		{
			description: "One empty log",
			testFunc: `
			func Test (t *testing.T) {
				t.Log("")
			}`,
			expected: []Test{},
		},

		{
			description: "One log",
			testFunc: `
			func Test (t *testing.T) {
				t.Log("something")
			}`,
			expected: []Test{},
		},

		{
			description: "Two logs",
			testFunc: `
			func Test (t *testing.T) {
				t.Log("something")
				t.Log("something else")
			}`,
			expected: []Test{},
		},

		{
			description: "Multiline log",
			testFunc: `
			func Test (t *testing.T) {
				t.Log("something\non\nfour\nlines")
			}`,
			expected: []Test{},
		},

		{
			description: "One empty fatal",
			testFunc: `
			func Test (t *testing.T) {
				t.Fatal("")
			}`,
			expected: []Test{
				{
					Name:   "Test",
					Lines:  []Line{{Number: "5", Output: ""}},
					Output: "",
				},
			},
		},

		{
			description: "One fatal",
			testFunc: `
			func Test (t *testing.T) {
				t.Fatal("error")
			}`,
			expected: []Test{
				{
					Name:   "Test",
					Lines:  []Line{{Number: "5", Output: "error"}},
					Output: "",
				},
			},
		},

		{
			description: "Two fatals",
			testFunc: `
			func Test (t *testing.T) {
				t.Fatal("error")
				t.Fatal("again")
			}`,
			expected: []Test{
				{
					Name:   "Test",
					Lines:  []Line{{Number: "5", Output: "error"}},
					Output: "",
				},
			},
		},

		{
			description: "Multiline fatal",
			testFunc: `
			func Test (t *testing.T) {
				t.Fatal("error\non\nfour\nlines")
			}`,
			expected: []Test{
				{
					Name:   "Test",
					Lines:  []Line{{Number: "5", Output: "error"}},
					Output: "\n        on\n        four\n        lines",
				},
			},
		},

		// {
		// // Any log printed before a failed test will be printed.
		// description: "Log then fatal",
		// testFunc: `
		// func Test (t *testing.T) {
		// t.Log("ok")
		// t.Fatal("error")
		// }`,
		// expected: []Test{
		// {Name: "Test", Line: "5", Output: "ok"},
		// {Name: "Test", Line: "6", Output: "error"},
		// },
		// },

		// {
		// description: "Log and fatal",
		// testFunc: `
		// func TestLog (t *testing.T) {
		// t.Log("ok")
		// }

		// func TestFatal (t *testing.T) {
		// t.Fatal("error")
		// }`,
		// expected: []Test{
		// {Name: "TestFatal", Line: "9", Output: "error"},
		// },
		// },

		// {
		// description: "Fatal and fatal",
		// testFunc: `
		// func TestFatal1 (t *testing.T) {
		// t.Fatal("error 1")
		// }

		// func TestFatal2 (t *testing.T) {
		// t.Fatal("error 2")
		// }`,
		// expected: []Test{
		// {Name: "TestFatal1", Line: "5", Output: "error 1"},
		// {Name: "TestFatal2", Line: "9", Output: "error 2"},
		// },
		// },
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			s, filename := createScanner(t, tc.testFunc)
			tests, err := Parse(s)

			assert.NoError(t, err, "parsing successful")

			assert.Equal(t, len(tc.expected), len(tests), "test count")
			t.Log(tests)

			for i := 0; i < len(tc.expected); i++ {
				assert.Equal(t, tc.expected[i].Output, tests[i].Output, "output")
				assert.Equal(t, filename, tests[i].File, "filename")
				assert.Equal(t, pkg, tests[i].Package, "package")

				for j, line := range tests[i].Lines {
					assert.Equal(t, tc.expected[i].Lines[j].Number, line.Number, "line number")
					assert.Equal(t, tc.expected[i].Lines[j].Output, line.Output, "line number")
				}
			}
		})
	}
}

func TestNotJSON(t *testing.T) {
	reader := strings.NewReader("not JSON")
	s := bufio.NewScanner(reader)
	tests, err := Parse(s)

	assert.Error(t, err)
	assert.Empty(t, tests)
}
