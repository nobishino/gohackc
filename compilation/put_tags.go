package compilation

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (e *Engine) indent() string {
	return strings.Repeat("  ", e.depth)
}

func (e *Engine) putNonTerminalTag(name string) func() {
	fmt.Println("OPEN putNonTerminalTag:", name)
	defer fmt.Println("CLOSE putNonTerminalTag:", name)

	noop := func() {}
	if _, err := io.WriteString(e.dst, e.indent()+"<"+name+">\n"); err != nil {
		e.addError(errors.WithStack(err))
		return noop
	}
	e.depth++
	return func() {
		e.depth--
		if _, err := io.WriteString(e.dst, e.indent()+"</"+name+">\n"); err != nil {
			e.addError(errors.WithStack(err))
		}
	}
}

func (e *Engine) putTerminalTag(name, value string) {
	if _, err := io.WriteString(e.dst, e.indent()+"<"+name+"> "+value+" </"+name+">\n"); err != nil {
		e.addError(errors.WithStack(err))
		return
	}
}

func (e *Engine) putKeywordTag(value string) {
	e.putTerminalTag("keyword", value)
}

func (e *Engine) putSymbolTag(value string) {
	switch value {
	case "<":
		e.putTerminalTag("symbol", "&lt;")
	case ">":
		e.putTerminalTag("symbol", "&gt;")
	case "&":
		e.putTerminalTag("symbol", "&amp;")
	default:
		e.putTerminalTag("symbol", value)
	}
}

func (e *Engine) putIdentifierTag(value string) {
	e.putTerminalTag("identifier", value)
}

func (e *Engine) putStringConstantTag(value string) {
	e.putTerminalTag("stringConstant", value)
}
func (e *Engine) putIntegerConstantTag(value int) {
	e.putTerminalTag("integerConstant", strconv.Itoa(value))
}
