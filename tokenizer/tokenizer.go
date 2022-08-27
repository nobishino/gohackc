package tokenizer

import (
	"bufio"
	"io"
	"strings"

	"github.com/nobishino/textfilter"
)

type Tokenizer struct {
	r             *bufio.Reader
	hasMoreTokens bool
	currentToken  token
}

// NewTokenizer
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r: bufio.NewReader(
			textfilter.NewReader(r, '\r'),
		),
		hasMoreTokens: true,
	}
}

func (t *Tokenizer) HasMoreTokens() bool {
	return t.hasMoreTokens
}

// 次のトークンを取得し、それをカレントトークンとする.
// HasMoreTokens()がtrueの場合のみ呼び出すことができる
func (t *Tokenizer) Advance() {
	r, _, err := t.r.ReadRune()
	if err == io.EOF {
		t.hasMoreTokens = false
		t.currentToken = token{tokenType: EOF}
		return
	}
	if err != nil {
		panic(err)
	}
	// skip spaces, tabs, and line separators
	for isDelimiter(r) {
		r, _, err = t.r.ReadRune()
		if err == io.EOF {
			t.hasMoreTokens = false
			t.currentToken = token{tokenType: EOF}
			return
		}
		if err != nil {
			panic(err)
		}
	}
	// symbol
	if isSymbol(r) {
		t.currentToken = token{
			tokenType: SYMBOL,
			symbol:    string(r),
		}
		return
	}
	// int_const
	// string_const
	// keyword
	// identifier
	panic("undefined")
}

func isSymbol(r rune) bool {
	return strings.ContainsRune("{}()[].,;+-*/&|<>=~", r)
}

func isDelimiter(r rune) bool {
	return strings.ContainsRune("\t\n ", r)
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
