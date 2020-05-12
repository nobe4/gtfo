package stuff_test

import "testing"

const (
	logMessage     = "log"
	fatalMessage   = "fatal"
	manyLinesLog   = "log \n on \n many \n lines"
	manyLinesFatal = "fatal \n on \n many \n lines"
)

func TestSuccessNoLog(t *testing.T) {
}

func TestSuccessOneLog(t *testing.T) {
	t.Log(logMessage)
}

func TestSuccessTwoLogs(t *testing.T) {
	t.Log(logMessage)
	t.Log(logMessage)
}

func TestSuccessMultilineLog(t *testing.T) {
	t.Log(manyLinesLog)
}

func TestFatalEmptyLog(t *testing.T) {
	t.Fatal("")
}

func TestFatalManyLines(t *testing.T) {
	t.Fatal(manyLinesFatal)
}

func TestManyFatalures(t *testing.T) {
	t.Fatal(fatalMessage)
	t.Fatal(fatalMessage)
}

func TestLogThenFatal(t *testing.T) {
	t.Log(logMessage)
	t.Fatal(fatalMessage)
}

func TestMultilineLogThenFatal(t *testing.T) {
	t.Log(manyLinesLog)
	t.Fatal(fatalMessage)
}

func TestMultilineLogThenMultilineFatal(t *testing.T) {
	t.Log(manyLinesLog)
	t.Fatal(manyLinesFatal)
}
