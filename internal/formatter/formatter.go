/*

Package formatter handlers how each logs will be displayed in the terminal.

Will enable user to configure the following keys to be replaced:

Key              | Name               | Example
---              | ---                | ---
{{.FullPackage}} | Absolute package   | github.com/nobe4/gtfo/internal/parser
{{.Package}}     | Local package      | internal/parser
{{.Module}}      | Project's module   | github.com/nobe4/gtfo
{{.File}}        | Filename           | parser_test.go
{{.FullPath}}    | Absolute file path | github.com/nobe4/gtfo/internal/parser/parser_test.go
{{.Path}}        | Relative file path | internal/parser/parser_test.go
{{.Line}}        | Log line           | 42
{{.Output}}      | Test output        | error\nsomething\nis\wrong

The replacement will use golang's template, read more here:
https://golang.org/pkg/text/template/

The format wll be passed to fmt.Print at the end, so all control characters
(e.g. `\n`, `\t`, ...) will be used normally. Make sure you include a `\n` at the end.

*/
package formatter

import (
	"bytes"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nobe4/gtfo/internal/parser"
)

type Fields struct {
	FullPackage string
	Package     string
	Module      string
	File        string
	FullPath    string
	Path        string
	Line        string
	Output      string
}

// escape replaces all the literral `\n`, `\t`, into their ASCII
// equivalent. I haven't found a way to pass a "\n" (i.e a newline) character
// in a consistent way, so we're getting `\n` (i.e. \ + n) and replacing them
// manually.
func escape(format string) string {
	escapeMap := map[string]string{
		`\n`:   "\n",
		`\r`:   "\r",
		`\r\n`: "\r\n",
		`\t`:   "\t",
	}

	for l, e := range escapeMap {
		format = strings.ReplaceAll(format, l, e)
	}

	return format
}

// Prepare creates a function to format the logs.
func Prepare(format, module string) (func(parser.Test) string, error) {
	format = escape(format)

	// Create the template
	t, err := template.New("test").Parse(format)
	if err != nil {
		return nil, err
	}

	return func(test parser.Test) string {
		fields := Fields{
			FullPackage: test.Package,
			Package:     strings.TrimPrefix(test.Package, module+"/"),
			Module:      module,
			File:        test.File,
			FullPath:    filepath.Join(test.Package, test.File),
			Line:        test.Line,
			Output:      test.Output,
		}

		fields.Path = filepath.Join(fields.Package, test.File)

		// Create a buffer to write into, so we can return the output of the
		// template.
		var buf bytes.Buffer

		// Execute the template
		err := t.Execute(&buf, fields)
		if err != nil {
			log.Fatal("Error executing template:", err)
		}

		return buf.String()
	}, nil
}
