package tokenizer

import (
	"fmt"
	"io"
	"strconv"
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
	if t.atEOF() {
		return 0
	}
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
	v, ok := t.readIntValue() // TODO: この書き方にまとめたほうがよいかもしれない
	if ok {
		t.currentToken = token{
			tokenType: INT_CONST,
			intValue:  v,
		}
		return
	}
	// string_const
	s, ok := t.readStringConstant()
	if ok {
		t.currentToken = token{
			tokenType:   STRING_CONST,
			stringValue: s,
		}
		return
	}
	// next token is keyword or identifier here
	w := t.readWord()
	// keyword
	if k, ok := isKeyword(w); ok {
		t.currentToken = token{
			tokenType: KEYWORD,
			keyword:   k,
		}
		return
	}
	// identifier
	t.currentToken = token{
		tokenType:  IDENTIFIER,
		identifier: w,
	}
}

func (t *Tokenizer) atEOF() bool {
	return t.pos >= len(t.sourceCode)
}

func isSymbol(r rune) bool {
	return strings.ContainsRune("{}()[].,;+-*/&|<>=~", r)
}

func (t *Tokenizer) skipDelimitersOrComments() {
	for t.atCommentStart() || t.atDelimiters() {
		t.skipBlockComment()
		t.skipLineComment()
		t.skipDelimiters()
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
		t.pos++
	}
	if t.atBlockCommentEnd() { // blockコメントが終わったときにソースコードの終端であるパターンに対応する FIXME: きれいにする
		return
	}
	panic("Block comment begins, but does not ends properly")
}

func isDelimiter(r rune) bool {
	return strings.ContainsRune("\r\t\n ", r)
}

// TokenType は、カレントトークンのタイプを返す
func (t *Tokenizer) TokenType() tokenType {
	return t.currentToken.tokenType
}

// 以下のメソッドは、 TokenType()の結果がそれぞれの前提とするTokenTypeであるときにのみよびだせる
// そうでないばあいはpanicする

// Keyword は、カレントトークンのKeyword値を返す
// TokenType()の値が不適切な場合はpanicする
func (t *Tokenizer) Keyword() KeyWord {
	if t.TokenType() != KEYWORD {
		msg := fmt.Sprintf("current token type must be %q but got %q", KEYWORD, t.TokenType())
		panic(msg)
	}
	return KeyWord(t.currentToken.keyword)
}

// Symbol は、カレントトークンのSymbol値を返す
// TokenType()の値が不適切な場合はpanicする
func (t *Tokenizer) Symbol() string {
	if t.TokenType() != SYMBOL {
		msg := fmt.Sprintf("current token type must be %q but got %q", SYMBOL, t.TokenType())
		panic(msg)
	}
	return t.currentToken.symbol
}

func (t *Tokenizer) Identifier() string {
	if t.TokenType() != IDENTIFIER {
		msg := fmt.Sprintf("current token type must be %q but got %q", IDENTIFIER, t.TokenType())
		panic(msg)
	}
	return t.currentToken.identifier
}

func (t *Tokenizer) IntVal() int {
	if t.TokenType() != INT_CONST {
		msg := fmt.Sprintf("current token type must be %q but got %q", INT_CONST, t.TokenType())
		panic(msg)
	}
	return t.currentToken.intValue
}

func (t *Tokenizer) StringVal() string {
	if t.TokenType() != STRING_CONST {
		msg := fmt.Sprintf("current token type must be %q but got %q", STRING_CONST, t.TokenType())
		panic(msg)
	}
	return t.currentToken.stringValue
}

// 第2戻り値はintegerConstantであるかどうかをかえす
// falseの場合はソースを読み進めない
// trueの場合はソースをintegerConstantの最後の文字まで読み終えてその値を返す
func (t *Tokenizer) readIntValue() (int, bool) {
	if !strings.ContainsRune("0123456789", t.currentLetter()) {
		return 0, false
	}
	begin := t.pos
	for !t.atEOF() && strings.ContainsRune("0123456789", t.currentLetter()) {
		t.pos++
	}
	v, err := strconv.Atoi(string(t.sourceCode[begin:t.pos]))
	if err != nil {
		panic(err)
	}
	return v, true
}

func (t *Tokenizer) readStringConstant() (string, bool) {
	if t.atEOF() || t.currentLetter() != '"' {
		return "", false
	}
	t.pos++
	begin := t.pos
	for !t.atEOF() {
		if t.currentLetter() == '"' {
			value := string(t.sourceCode[begin:t.pos])
			t.pos++
			return value, true
		}
		t.pos++
	}
	panic(`string constant does not closed:` + t.String())
}

func (t *Tokenizer) readWord() string {
	if t.atEOF() {
		panic("cannot happen")
	}
	begin := t.pos
	for !t.atEOF() && !t.atDelimiters() && !isSymbol(t.currentLetter()) {
		t.pos++
	}
	w := t.sourceCode[begin:t.pos]
	return w
}

func (t *Tokenizer) String() string {
	return fmt.Sprintf(`{
	sourceCode:	%q,
	cursor position:	%d/%d,
	current token:	%v,
	current letter: %c,
}`,
		t.sourceCode, t.pos, len(t.sourceCode), t.currentToken, t.currentLetter())
}
