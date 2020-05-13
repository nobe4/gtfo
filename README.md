GTFO
====

> Go Test FOrmat

This package aims at creating a simple formatter for the test logs in golang.

In a sense, this adds a missing `-format` to the `go test` command. This is
required in some cases when the editor expect some defined-output from a
command in order to jump directly to the issues.

# Usage

`gtfo` accepts only the JSON output from `go test`, so you can call it in the following way:

```
go test -json ./... | gtfo
```

All the flags for `test` should work as well, tho you might not get the
expected output. I recommend using `gtfo` for testing and jumping to results
directly, but not for complete output.

# Example
If we make a test in
[`parser_test.go`](./internal/parser/parser_test.go) fail, we might get the following output:

```bash
$ go test ./...
?       github.com/nobe4/gtfo/cmd/gtfo  [no test files]
ok      github.com/nobe4/gtfo/internal/module   (cached)
--- FAIL: TestCases (17.26s)
    parser_test.go:185: No Test
    ...
    parser_test.go:185: Fatal and fatal
    parser_test.go:194:
                Error Trace:    parser_test.go:194
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
FAIL    github.com/nobe4/gtfo/internal/parser   17.445s
FAIL
```

```bash
$ go test -json ./... | gtfo
internal/module/parser_test.go:185: No Test
...
internal/module/parser_test.go:185: Fatal and fatal
internal/module/parser_test.go:194:
                Error Trace:    parser_test.go:194
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
