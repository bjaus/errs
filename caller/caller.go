package caller

import (
	"fmt"
	"runtime"
	"strings"
)

type Caller interface {
	Package() string
	Function() string
	LineNumber() int
	String() string
}

var _ Caller = new(caller)

type caller struct {
	pkg  string
	fn   string
	path string
	line int
}

func (c caller) Package() string {
	return c.pkg
}

func (c caller) Function() string {
	return c.fn
}

func (c caller) FilePath() string {
	return c.path
}

func (c caller) LineNumber() int {
	return c.line
}

func (c caller) String() string {
	return fmt.Sprintf("%s.%s:%d", c.pkg, c.fn, c.line)
}

var replacer = *strings.NewReplacer("(*", "", ")", "", ".go", "")

func Parse(skip int) caller {
	c := caller{
		pkg:  "unknown",
		fn:   "unknown",
		line: -1,
		path: "unknown",
	}

	if skip < 0 {
		skip = 0
	}
	skip++

	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return c
	}

	f := runtime.FuncForPC(pc)
	parts := strings.Split(f.Name(), "/")

	caller := replacer.Replace(parts[len(parts)-1])
	parts = strings.Split(caller, ".func")

	caller = parts[0]
	parts = strings.Split(caller, ".")

	switch len(parts) {
	case 0:
		return c
	case 1:
		c.pkg = parts[0]
	default:
		c.pkg = parts[0]
		c.fn = strings.Join(parts[1:], ".")
	}

	c.line = line
	c.path = file

	return c
}
