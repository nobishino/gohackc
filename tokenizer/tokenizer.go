package tokenizer

import (
	"io"
	"strings"
)

type Tokenizer struct {
	hasMoreTokens bool
	currentToken  token
	sourceCode    string
	pos           int
}

// NewTokenizer
func NewTokenizer(r io.Reader) *Tokenizer {
	var sb strings.Builder
	if _, err := io.Copy(&sb, r); err != nil {
		panic(err)
	}
	return &Tokenizer{
		hasMoreTokens: true,
		sourceCode:    sb.String(),
	}
}

func (t *Tokenizer) HasMoreTokens() bool {
	return t.hasMoreTokens
}

// 次のトークンを取得し、それをカレントトークンとする.
// HasMoreTokens()がtrueの場合のみ呼び出すことができる
func (t *Tokenizer) Advance() {
	t.skipDelimiters()
	if t.pos >= len(t.sourceCode) {
		t.hasMoreTokens = false
		t.currentToken = token{tokenType: EOF}
		return
	}
	// if err != nil {
	// 	panic(err)
	// }
	// skip spaces, tabs, and line separators
	// symbol

	// int_const
	// string_const
	// keyword
	// identifier
	panic("undefined")
}

func isSymbol(r rune) bool {
	return strings.ContainsRune("{}()[].,;+-*/&|<>=~", r)
}

func (t *Tokenizer) skipDelimiters() {
	for t.pos < len(t.sourceCode) && isDelimiter(rune(t.sourceCode[t.pos])) {
		t.pos++
	}
}

func isDelimiter(r rune) bool {
	return strings.ContainsRune("\r\t\n ", r)
}

// TokenType
func (t *Tokenizer) TokenType() tokenType {
	panic("undefined")
}

// 以下のメソッドは、 TokenType()の結果がそれぞれの前提とするTokenTypeであるときにのみよびだせる
// そうでないばあいはpanicする

// Keyword は、カレントトークンのKeyword値を返す
// TokenType()の値が不適切な場合はpanicする
func (t *Tokenizer) Keyword() keyWord {
	panic("undefined")
}

// Symbol は、カレントトークンのSymbol値を返す
// TokenType()の値が不適切な場合はpanicする
func (t *Tokenizer) Symbol() string {
	panic("undefined")
}

func (t *Tokenizer) Identifier() string {
	panic("undefined")
}

func (t *Tokenizer) IntVal() int {
	panic("undefined")
}

func (t *Tokenizer) StringVal() string {
	panic("undefined")
}
