package tokenizer

import "io"

type Tokenizer struct{}

// NewTokenizer
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{}
}

func (t *Tokenizer) HasMoreTokens() bool {
	panic("undefined")
}

// 次のトークンを取得し、それをカレントトークンとする.
// HasMoreTokens()がtrueの場合のみ呼び出すことができる
func (t *Tokenizer) Advance() {
	panic("undefined")
}

// TokenType
func (t *Tokenizer) TokenType() TokenType {
	panic("undefined")
}

// 以下のメソッドは、 TokenType()の結果がそれぞれの前提とするTokenTypeであるときにのみよびだせる
// そうでないばあいはpanicする

// Keyword は、カレントトークンのKeyword値を返す
// TokenType()の値が不適切な場合はpanicする
func (t *Tokenizer) Keyword() Keyword {
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
