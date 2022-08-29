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

func (t *Tokenizer) currentLetter() rune {
	return rune(t.sourceCode[t.pos])
}

// 次のトークンを取得し、それをカレントトークンとする.
// HasMoreTokens()がtrueの場合のみ呼び出すことができる
func (t *Tokenizer) Advance() {
	if t.atEOF() {
		t.hasMoreTokens = false
		t.currentToken = token{tokenType: EOF}
		return
	}
	// skip spaces, tabs, and line separators
	t.skipDelimitersOrComments()
	if t.atEOF() {
		t.hasMoreTokens = false
		t.currentToken = token{tokenType: EOF}
		return
	}
	// symbol
	if isSymbol(t.currentLetter()) {
		t.currentToken = token{
			tokenType: SYMBOL,
			symbol:    string(t.currentLetter()),
		}
		t.pos++
		return
	}
	// int_const
	// string_const
	// keyword
	// identifier
	panic("undefined")
}

func (t *Tokenizer) atEOF() bool {
	return t.pos >= len(t.sourceCode)
}

func isSymbol(r rune) bool {
	return strings.ContainsRune("{}()[].,;+-*/&|<>=~", r)
}

func (t *Tokenizer) skipDelimitersOrComments() {
	for t.atCommentStart() || t.atDelimiters() {
		t.skipDelimiters()
		t.skipLineComment()
		t.skipBlockComment()
	}
}

func (t *Tokenizer) skipDelimiters() {
	for t.pos < len(t.sourceCode) && isDelimiter(rune(t.sourceCode[t.pos])) {
		t.pos++
	}
}

func (t *Tokenizer) atDelimiters() bool {
	return !t.atEOF() && isDelimiter(t.currentLetter())
}

func (t *Tokenizer) skipLineComment() {
	if t.atCommentStart() {
		t.skipUntilLF()
	}
}

// いまいる位置にコメントの開始シーケンスが存在するかどうかを返す
func (t *Tokenizer) atCommentStart() bool {
	if t.atLineCommentStart() {
		return true
	}
	return t.atBlockCommentStart()
}

func (t *Tokenizer) atLineCommentStart() bool {
	if t.pos >= len(t.sourceCode)-1 { // 最後の1文字ならコメントの開始にはならない
		return false
	}
	return t.sourceCode[t.pos:t.pos+2] == "//"
}

func (t *Tokenizer) skipBlockComment() {
	if t.atBlockCommentStart() {
		t.skipUntilBlockCommentEnd()
	}
}

func (t *Tokenizer) atBlockCommentStart() bool {
	if t.pos >= len(t.sourceCode)-2 { // 最後の2文字以降ならブロックコメントの開始にはならない
		return false
	}
	return t.sourceCode[t.pos:t.pos+3] == "/**"
}

func (t *Tokenizer) atBlockCommentEnd() bool {
	if t.pos < 2 { // 最低でも3文字目以降でなければならない
		return false
	}
	return t.sourceCode[t.pos-2:t.pos] == "*/"
}

func (t *Tokenizer) skipUntilLF() {
	for t.pos < len(t.sourceCode) && t.currentLetter() != '\n' {
		t.pos++
	}
}

func (t *Tokenizer) skipUntilBlockCommentEnd() {
	for t.pos < len(t.sourceCode) {
		if t.atBlockCommentEnd() {
			return
		}
	}
	panic("Block comment begins, but does not ends properly")
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
