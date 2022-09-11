// Package compilation はコンパイルを行う
package compilation

import (
	"fmt"
	"io"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/nobishino/gohackc/tokenizer"
	"github.com/pkg/errors"
)

type Engine struct {
	tz    *tokenizer.Tokenizer
	dst   io.Writer
	errs  error
	depth int
}

func New(src io.Reader, dst io.Writer) *Engine {
	tz := tokenizer.NewTokenizer(src)
	e := &Engine{
		tz:  tz,
		dst: dst,
	}
	if e.tz.HasMoreTokens() {
		e.advance()
	}
	return e
}

// class = 'class' className '{' classVarDec* subroutineDec* '}'
func (e *Engine) CompileClass() {
	if ok := e.eat("class"); !ok {
		e.addError(errors.Errorf("error: expect keyword %q but currnt token is not", "class"))
		return
	}
	closer := e.putNonTerminalTag("class")
	defer closer()

	e.putKeywordTag("class")

	if e.tz.TokenType() != tokenizer.IDENTIFIER {
		e.addError(errors.Errorf("expect identifier as className, but got %q", e.tz.TokenType()))
		return
	}
	e.putIdentifierTag(e.tz.Identifier())
	e.advance()

	if ok := e.eat("{"); !ok {
		e.addError(errors.Errorf("error: expect symbol %q but currnt token is not", "{"))
		return
	}
	e.putSymbolTag("{")

	// TODO: classの中身をparseする
	if ok := e.eat("}"); !ok {
		e.addError(errors.Errorf("error: expect symbol %q but currnt token is not", "}"))
		return
	}
	e.putSymbolTag("}")
}

// classVarDec = ('static' | 'field' ) type varName (',' varName)* ';'
func (e *Engine) compileClassVarDec() error {
	return nil
}

// subroutineDec = ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
func (e *Engine) compileSubroutine() error {
	return nil
}

func (e *Engine) compileParameterList() error {
	return nil
}

func (e *Engine) compileVarDec() error {
	return nil
}

func (e *Engine) compileStatements() error {
	return nil
}

func (e *Engine) compileDo() error {
	return nil
}

func (e *Engine) compileLet() error {
	return nil
}

func (e *Engine) compileWhile() error {
	return nil
}

func (e *Engine) compileReturn() error {
	return nil
}

func (e *Engine) compileIf() error {
	return nil
}

func (e *Engine) compileExpression() error {
	return nil
}

func (e *Engine) compileTerm() error {
	return nil
}

func (e *Engine) compileExpressionList() error {
	return nil
}

func (e *Engine) eat(value string) bool {
	switch kind := e.tz.TokenType(); kind {
	case tokenizer.KEYWORD:
		if e.tz.Keyword() == tokenizer.KeyWord(value) {
			if e.tz.HasMoreTokens() {
				e.advance()
			}
			return true
		}
	case tokenizer.SYMBOL:
		if e.tz.Symbol() == value {
			if e.tz.HasMoreTokens() {
				e.advance()
			}
		}
		return true
	case tokenizer.IDENTIFIER, tokenizer.INT_CONST, tokenizer.STRING_CONST, tokenizer.EOF:
		e.addError(errors.Errorf("eat cannot used with token type %q", kind))
	default:
		e.addError(errors.Errorf("eat cannot used with token type %q", kind))
	}
	return false
}

func (e *Engine) addError(err error) {
	e.errs = multierror.Append(e.errs, err)
}

func (e *Engine) Error() error {
	return e.errs
}

func (e *Engine) advance() {
	e.tz.Advance()
	e.logCurrentToken()
}

func (e *Engine) logCurrentToken() {
	prefix := "[logCurrentToken] "
	switch tp := e.tz.TokenType(); tp {
	case tokenizer.KEYWORD:
		log.Printf(prefix+"Keyword: %q", e.tz.Keyword())
	case tokenizer.SYMBOL:
		log.Printf(prefix+"Symbol: %q", e.tz.Symbol())
	case tokenizer.IDENTIFIER:
		log.Printf(prefix+"Identifier: %q", e.tz.Identifier())
	case tokenizer.INT_CONST:
		log.Printf(prefix+"Integer Constant: %q", e.tz.Identifier())
	case tokenizer.STRING_CONST:
		log.Printf(prefix+"String constant: %q", e.tz.Identifier())
	case tokenizer.EOF:
		log.Print(prefix + "EOF")
	default:
		panic(fmt.Sprintf(prefix+"unexpected token type %q", tp))
	}
}
