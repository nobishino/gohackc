// Package compilation はコンパイルを行う
package compilation

import (
	"io"

	"github.com/nobishino/gohackc/tokenizer"
)

type Engine struct {
	tz  *tokenizer.Tokenizer
	dst io.Writer
}

func New(src io.Reader, dst io.Writer) *Engine {
	tz := tokenizer.NewTokenizer(src)
	return &Engine{
		tz:  tz,
		dst: dst,
	}
}

// class = 'class' className '{' classVarDec* subroutineDec* '}'
func (e *Engine) CompileClass() error {
	return nil
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
