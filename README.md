GTFO
====

> Go Test FOrmat

This package aims at creating a simple formatter for the test logs in golang.

In a sense, this adds a missing `-format` to the `go test` command. This is
required in some cases when the editor expect some defined-output from a
command in order to jump directly to the issues.

# Usage

`gtfo` accepts only the JSON output from `go test`, so you can call it in the
following way:

```
go test -json ./... | gtfo
```

All the flags for `test` should work as well, tho you might not get the
expected output. I recommend using `gtfo` for testing and jumping to results
directly, but not for complete output.

# Configuration

You can pass a string to the `-format` flag to specify how you want to format
your output. The syntax follows [golang's
template](https://golang.org/pkg/text/template/). Have a look at
[formatter.go](./internal/formatter/formatter.go) for more informations.

# Example
If we make a test in
[`parser_test.go`](./internal/parser/parser_test.go) fail, we might get the
following output:

```bash
$ go test ./...
?       github.com/nobe4/gtfo/cmd/gtfo  [no test files]
ok      github.com/nobe4/gtfo/internal/formatter        0.523s
ok      github.com/nobe4/gtfo/internal/module   (cached)
--- FAIL: TestCases (14.80s)
    parser_test.go:188: No Test
    ...
    parser_test.go:188: Fatal and fatal
    parser_test.go:197:
                Error Trace:    parser_test.go:197
                Error:          Not equal:
                                expected: "8"
                                actual  : "9"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -8
                                +9
                Test:           TestCases
FAIL
FAIL    github.com/nobe4/gtfo/internal/parser   15.223s
FAIL
```

```bash
$ go test -json ./... | gtfo
internal/formatter/parser_test.go:188: No Test
...
internal/formatter/parser_test.go:188: Fatal and fatal
internal/formatter/parser_test.go:197:
                Error Trace:    parser_test.go:197
                Error:          Not equal:
                                expected: "8"
                                actual  : "9"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -8
                                +9
                Test:           TestCases
```

We have the full path from the root directory to the file in question, which
allows for easier jumping to it.

With a special format:

```bash
$ go test -json ./... | gtfo -format "{{.File}}:{{.Line}}:{{.Output}}\n"
parser_test.go:188:No Test
...
parser_test.go:188:Fatal and fatal
parser_test.go:197:
                Error Trace:    parser_test.go:197
                Error:          Not equal:
                                expected: "8"
                                actual  : "9"

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1 +1 @@
                                -8
                                +9
                Test:           TestCases
```



# Installation

Clone and build:
```bash
make build
```

You might want to change some environment variables, e.g.:
```bash
GOOS=linux make build
```
# Contribute

Ideas and improvement in [TODO](./TODO.md).
