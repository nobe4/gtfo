package formatter

import (
	"testing"

	"github.com/nobe4/gtfo/internal/parser"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestFormat(t *testing.T) {
	testCases := []struct {
		description string
		format      string
		module      string
		test        parser.Test
		expected    string
	}{
		{
			description: "No format",
			// expected == ""
		},

		{
			description: "Fixed format",
			format:      "test",
			expected:    "test",
		},

		{
			description: "Full package format",
			format:      "{{.FullPackage}}",
			test:        parser.Test{Package: "package"},
			expected:    "package",
		},

		{
			description: "Package format",
			module:      "module",
			format:      "{{.Package}}",
			test:        parser.Test{Package: "module/package"},
			expected:    "package",
		},

		{
			description: "Full path format",
			format:      "{{.FullPath}}",
			test:        parser.Test{Package: "module/package", File: "file"},
			expected:    "module/package/file",
		},

		{
			description: "Path format",
			module:      "module",
			format:      "{{.Path}}",
			test:        parser.Test{Package: "module/package", File: "file"},
			expected:    "package/file",
		},

		{
			description: "All format",
			format:      "{{.FullPackage}} {{.Package}} {{.Module}} {{.File}} {{.FullPath}} {{.Path}} {{.Line}} {{.Output}}",
			test: parser.Test{
				Package: "module/package",
				File:    "file",
				Line:    "1",
				Output:  "output",
			},
			module:   "module",
			expected: "module/package package module file module/package/file package/file 1 output",
		},

		{
			description: "Formatting",
			format:      "{{.Module}}\n{{.Module}}\t{{.Module}}",
			module:      "module",
			expected:    "module\nmodule\tmodule",
		},
	}
	for _, tc := range testCases {
		t.Log(tc.description)
		format, err := Prepare(tc.format, tc.module)
		assert.NoError(t, err)

		result := format(tc.test)
		assert.Equal(t, tc.expected, result)
	}
}

func TestErrorTemplate(t *testing.T) {
	format, err := Prepare("{{.Missing End }", "")
	assert.Error(t, err)
	assert.Nil(t, format)
}

func TestEscape(t *testing.T) {
	testCases := []struct {
		description string
		format      string
		expected    string
	}{
		{
			description: "Empty",
		},
		{
			description: "No special chars",
			format:      "abcd",
			expected:    "abcd",
		},
		{
			description: "Newline",
			format:      `a\nb`,
			expected:    "a\nb",
		},
		{
			description: "Newline and tab",
			format:      `a\nb\tc`,
			expected:    "a\nb\tc",
		},
	}

	for _, tc := range testCases {
		t.Log(tc.description)
		format := escape(tc.format)
		assert.Equal(t, tc.expected, format)
	}
}
