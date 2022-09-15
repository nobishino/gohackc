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

ParseClassVarDecs:
	for e.tz.TokenType() == tokenizer.KEYWORD {
		switch e.tz.Keyword() {
		case "static", "field":
			e.compileClassVarDec()
		default:
			break ParseClassVarDecs
		}
	}

	for e.tz.TokenType() == tokenizer.KEYWORD {
		switch e.tz.Keyword() {
		case "constructor", "function", "method":
			e.compileSubroutineDec()
		default:
			e.addError(errors.Errorf("unexpected keyword: %q", e.tz.Keyword()))
			return
		}
	}

	if ok := e.eat("}"); !ok {
		e.addError(errors.Errorf("error: expect symbol %q but currnt token is not", "}"))
		return
	}
	e.putSymbolTag("}")
}

// classVarDec = ('static' | 'field' ) type varName (',' varName)* ';'
func (e *Engine) compileClassVarDec() {
	defer e.putNonTerminalTag("classVarDec")()
	kw := e.tz.Keyword()
	e.putKeywordTag(string(kw))
	e.advance()

	// typeはKeywordとIdentifierどちらかがくる
	switch tokenType := e.tz.TokenType(); tokenType {
	case tokenizer.KEYWORD:
		classVarType, ok := e.expectKeyword()
		if !ok {
			panic("its bug")
		}
		e.putKeywordTag(classVarType)
	case tokenizer.IDENTIFIER:
		classVarType, ok := e.expectIdentifier()
		if !ok {
			panic("its bug")
		}
		e.putIdentifierTag(classVarType)
	default:
		e.addError(errors.Errorf("unexpected type: %q", tokenType))
		return
	}

ParseVarNames:
	for {
		ident, ok := e.expectIdentifier()
		if !ok {
			return
		}
		e.putIdentifierTag(ident)

		if e.tz.TokenType() != tokenizer.SYMBOL {
			e.addError(errors.Errorf("expect symbol (; or ,), but got type %q", e.tz.TokenType()))
			return
		}

		switch s := e.tz.Symbol(); s {
		case ";":
			e.advance()
			e.putSymbolTag(s)
			break ParseVarNames
		case ",":
			e.advance()
		default:
			e.addError(errors.Errorf("unexpected symbol value: %q", s))
			return
		}
	}
}

// subroutineDec = ('constructor' | 'function' | 'method') ('void' | type) subroutineName '(' parameterList ')' subroutineBody
func (e *Engine) compileSubroutineDec() {
	closer := e.putNonTerminalTag("subroutineDec")

	// 'constructor' | 'function' | 'method'
	kw, ok := e.expectKeyword()
	if !ok {
		return
	}
	e.putKeywordTag(kw)

	// 'void' | 'type'
	switch tokenType := e.tz.TokenType(); tokenType {
	case tokenizer.KEYWORD:
		kw, _ := e.expectKeyword()
		switch kw {
		case "void", "int", "char", "boolean":
			e.putKeywordTag(kw)
		default:
			e.addError(errors.Errorf("unexpected keyword: %q", kw))
		}
	case tokenizer.IDENTIFIER:
		className, _ := e.expectIdentifier()
		e.putIdentifierTag(className)
	default:
		e.addError(errors.Errorf("unexpected token type %q, expecting ('void' | type).", tokenType))
	}

	// subroutineName
	subroutineName, ok := e.expectIdentifier()
	if !ok {
		return
	}
	e.putIdentifierTag(subroutineName)

	// '(' parameterList ')'
	if !e.eatSymbol("(") {
		return
	}
	e.putSymbolTag("(")
	e.compileParameterList()
	if !e.eatSymbol(")") {
		return
	}
	e.putSymbolTag(")")

	closeSubroutineBody := e.putNonTerminalTag("subroutineBody")

	if !e.eatSymbol("{") {
		return
	}
	e.putSymbolTag("{")

	// varDec*
	for e.tz.TokenType() == tokenizer.KEYWORD && e.tz.Keyword() == "var" {
		e.compileVarDec()
	}

	// statements
	e.compileStatements()

	if !e.eatSymbol("}") {
		return
	}
	e.putSymbolTag("}")

	closeSubroutineBody()

	closer()
}

// parameterList = ((type varName) (',' type VarName)*)?
// type = 'int' | 'char' | 'boolean' | className
// className = identifier
// varName = identifier
func (e *Engine) compileParameterList() {
	defer e.putNonTerminalTag("parameterList")()
	// parameterListが空であることとこの時点でSymbolであることは同値
	if e.tz.TokenType() == tokenizer.SYMBOL {
		return
	}
	// ここに来る場合は1つ以上の(type varName)の繰り返しになる
parameters:
	for {
		switch tokenType := e.tz.TokenType(); tokenType {
		case tokenizer.KEYWORD:
			kw, _ := e.expectKeyword()
			switch kw {
			case "int", "boolean", "char":
				e.putKeywordTag(kw)
			default:
				e.addError(errors.Errorf("unexpected keyword: %q", kw))
				return
			}
		case tokenizer.IDENTIFIER:
			typeName, _ := e.expectIdentifier()
			e.putIdentifierTag(typeName)
		}
		varName, ok := e.expectIdentifier()
		if !ok {
			return
		}
		e.putIdentifierTag(varName)

		// type VarNameが2個以上のパターンに対応する
		if e.tz.TokenType() == tokenizer.SYMBOL && e.tz.Symbol() != "," {
			break parameters
		} else {
			e.eatSymbol(",")
			e.putSymbolTag(",")
		}
	}
}

// varDec = 'var' type identifier;
func (e *Engine) compileVarDec() {
	varDecCloser := e.putNonTerminalTag("varDec")
	defer varDecCloser()
	// 'var'
	if !e.eatKeyword("var") {
		return
	}
	e.putKeywordTag("var")

	// type
	e.helpCompileType()

	// identifer
	ident, ok := e.expectIdentifier()
	if !ok {
		return
	}
	e.putIdentifierTag(ident)

	// ;
	if !e.eat(";") {
		return
	}
	e.putSymbolTag(";")
}

// type =
func (e *Engine) helpCompileType() {
	switch tp := e.tz.TokenType(); tp {
	case tokenizer.IDENTIFIER:
		typeName, _ := e.expectIdentifier()
		e.putIdentifierTag(typeName)
	case tokenizer.KEYWORD:
		switch kw, _ := e.expectKeyword(); kw {
		case "int", "boolean", "char":
			e.putKeywordTag(kw)
		default:
			e.addError(errors.Errorf("expect int, boolean, or char. but got keyword %q", kw))
		}
	}
}

func (e *Engine) compileStatements() {
	defer e.putNonTerminalTag("statements")()
Statements:
	for e.tz.TokenType() == tokenizer.KEYWORD {
		switch e.tz.Keyword() {
		case tokenizer.RETURN:
			e.compileReturn()
		case tokenizer.WHILE:
			e.compileWhile()
		case tokenizer.LET:
			e.compileLet()
		case tokenizer.IF:
			e.compileIf()
		case tokenizer.DO:
			e.compileDo()
		default:
			break Statements
		}
	}
}

// 'do' subroutineCall
func (e *Engine) compileDo() {
	closeDo := e.putNonTerminalTag("doStatement")
	defer closeDo()

	// do
	if !e.eatKeyword("do") {
		return
	}
	e.putKeywordTag("do")

	// subroutineCall
	e.helpCompileSubroutineCall()

	// ;
	if !e.eat(";") {
		return
	}
	e.putSymbolTag(";")
}

// subroutineCall = subroutineName '(' expressionList ')' | (className | varName) '.' subroutineName '(' expressionList ')'
// className, varName, subroutineNameはすべてidentifier
func (e *Engine) helpCompileSubroutineCall() {
	// 開始終了タグはなし

	// subroutineName, className, varName
	ident, ok := e.expectIdentifier()
	if !ok {
		return
	}
	e.putIdentifierTag(ident)

	// must be '.' or '('
	s, ok := e.expectSymbol()
	if !ok {
		return
	}
	switch s {
	case ".": // (className | varName) . subroutineName ( expressionList )
		e.putSymbolTag(s)
		subRoutineName, ok := e.expectIdentifier()
		if !ok {
			return
		}
		e.putIdentifierTag(subRoutineName)
		if !e.eat("(") {
			return
		}
		e.putSymbolTag("(")
		e.compileExpressionList()
		if !e.eat(")") {
			return
		}
		e.putSymbolTag(")")
	case "(": // subroutineName '(' expressionList ')'
		e.putSymbolTag(s)
		e.compileExpressionList()
		if !e.eat(")") {
			return
		}
		e.putSymbolTag(")")
	default:
		e.addError(errors.Errorf("subroutineCall: should be '(' or '.' here, but got %q", s))
		return
	}

}

// 'let' varName ('[' expression ']')? '=' expression;
// varName = identifier
func (e *Engine) compileLet() {
	closeLetStatement := e.putNonTerminalTag("letStatement")
	defer closeLetStatement()
	// let
	if !e.eat("let") {
		return
	}
	e.putKeywordTag("let")
	// varName (左辺)
	varName, ok := e.expectIdentifier()
	if !ok {
		return
	}
	e.putIdentifierTag(varName)
	// [ がなければindexingではない
	if e.peekKeyword("[") {
		panic("not implemented")
	}
	// =
	if !e.eat("=") {
		return
	}
	e.putSymbolTag("=")
	// expression (右辺)
	e.compileExpression()
	if !e.eat(";") {
		return
	}
	e.putSymbolTag(";")
}

// whileStatement = 'while' '(' expression ')' '{' statements '}'
func (e *Engine) compileWhile() {
	closeWhile := e.putNonTerminalTag("whileStatement")
	defer closeWhile()

	// while
	if !e.eatKeyword("while") {
		return
	}
	e.putKeywordTag("while")

	// ( expression )
	if !e.eatSymbol("(") {
		return
	}
	e.putSymbolTag("(")
	e.compileExpression()
	if !e.eatSymbol(")") {
		return
	}
	e.putSymbolTag(")")

	if !e.eatSymbol("{") {
		return
	}
	e.putSymbolTag("{")
	e.compileStatements()
	if !e.eatSymbol("}") {
		return
	}
	e.putSymbolTag("}")
}

func (e *Engine) compileReturn() {
	defer e.putNonTerminalTag("returnStatement")()
	if !e.eatKeyword("return") {
		return
	}
	e.putKeywordTag("return")
	if !(e.tz.TokenType() == tokenizer.SYMBOL && e.tz.Symbol() == ";") {
		e.compileExpression()
	}
	if !e.eatSymbol(";") {
		return
	}
	e.putSymbolTag(";")
}

// if '(' expression ')' '{' statements '}' ( 'else' '{' statements '}' )?
func (e *Engine) compileIf() {
	closeIf := e.putNonTerminalTag("ifStatement")
	defer closeIf()
	// if
	if !e.eatKeyword("if") {
		return
	}
	e.putKeywordTag("if")

	// ( expression )
	if !e.eatSymbol("(") {
		return
	}
	e.putSymbolTag("(")
	e.compileExpression()
	if !e.eatSymbol(")") {
		return
	}
	e.putSymbolTag(")")

	// { statements }
	e.helpCompileStatementsWithCurlyBrackets()

	// (else { statements })?
	if e.tz.TokenType() == tokenizer.KEYWORD && e.tz.Keyword() == "else" {
		e.eatKeyword("else")
		e.putKeywordTag("else")
		e.helpCompileStatementsWithCurlyBrackets()
	}

}

func (e *Engine) helpCompileStatementsWithCurlyBrackets() {
	if !e.eatSymbol("{") {
		return
	}
	e.putSymbolTag("{")
	e.compileStatements()
	if !e.eatSymbol("}") {
		return
	}
	e.putSymbolTag("}")
}

// expression = term (op term)*
// op = +, -, *, /, &, |, <, >, =
func (e *Engine) compileExpression() {
	defer e.putNonTerminalTag("expression")()
	e.compileTerm()
	for {
		op, ok := e.maybeOp()
		if !ok {
			break
		}
		e.putSymbolTag(op)
		e.compileTerm()
	}
}

// term = intergerConstant
//
//	| stringConstant
//	| keywordConstant
//	| varName
//	| varName '[' expression ']'
//	| subroutineCall
//	| '(' expression ')'
//	| unaryOp term
//
// unaryOp = '-' | '~'
func (e *Engine) compileTerm() {
	closeTerm := e.putNonTerminalTag("term")
	defer closeTerm()
	switch tokenType := e.tz.TokenType(); tokenType {
	// integerConstant
	case tokenizer.INT_CONST:
		intConst, _ := e.expectIntegerConstant()
		e.putIntegerConstantTag(intConst)
	// stringConstant
	case tokenizer.STRING_CONST:
		stringConst, _ := e.expectIntegerConstant()
		e.putIntegerConstantTag(stringConst)
	// keywordConstant
	case tokenizer.KEYWORD:
		switch kw, _ := e.expectKeyword(); kw {
		case "true", "false", "null", "this":
			e.putKeywordTag(kw)
		default:
			e.addError(errors.Errorf("expect keyword constant but got keyword %q", kw))
			return
		}
	case tokenizer.SYMBOL:
		// unaryOp term
		switch symbol, _ := e.expectSymbol(); symbol {
		case "-", "~":
			e.putSymbolTag(symbol)
		// '(' expression ')'
		case "(":
			e.putSymbolTag(symbol)
			e.compileExpression()
			e.eatSymbol(")")
			e.putSymbolTag(")")
		}
	// varName
	case tokenizer.IDENTIFIER:
		varName, _ := e.expectIdentifier()
		e.putIdentifierTag(varName)
		// TODO: varName '[' expression ']'
		// TODO: subroutineCall
	}
}

// expressionList = (expression (',' expression)* )?
// これむずかしいのでは
// TODO: implement
// 現在は空のパターンのみサポート
func (e *Engine) compileExpressionList() {
	closeExpressionList := e.putNonTerminalTag("expressionList")
	defer closeExpressionList()
	// FIXME:
	// expressionが１つもない場合はsymbol ) がカレントトークンであると仮定する←ただしい？
	if e.tz.TokenType() == tokenizer.SYMBOL && e.tz.Symbol() == ")" {
		return
	}
	for {
		e.compileExpression()
		if e.tz.TokenType() == tokenizer.SYMBOL && e.tz.Symbol() == "," {
			e.advance()
		} else {
			return
		}
	}
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

// カレントトークンがvalueという値のKeywordである場合はtrueを返して次のトークンに進む
// そうでない場合はエラーを追加してfalseを返す
func (e *Engine) eatKeyword(value string) bool {
	if e.tz.TokenType() != tokenizer.KEYWORD {
		e.addError(errors.Errorf("eatKeyword() was given token type %q", e.tz.TokenType()))
		return false
	}
	if e.tz.Keyword() != tokenizer.KeyWord(value) {
		e.addError(errors.Errorf("eatKeyword() expect %q, but got %q", value, e.tz.Keyword()))
		return false
	}
	e.advance()
	return true
}

// カレントトークンがKeywordである場合はその値とtrueを返して次のトークンに進む
// そうでない場合はエラーを追加してfalseを返す
func (e *Engine) expectKeyword() (string, bool) {
	if e.tz.TokenType() != tokenizer.KEYWORD {
		e.addError(errors.Errorf("expectKeyword() was given token type %q", e.tz.TokenType()))
		return "", false
	}
	value := string(e.tz.Keyword())
	e.advance()
	return value, true
}

func (e *Engine) peekKeyword(kw string) bool {
	if e.tz.TokenType() != tokenizer.KEYWORD {
		return false
	}
	return string(e.tz.Keyword()) == kw
}

func (e *Engine) eatSymbol(value string) bool {
	if e.tz.TokenType() != tokenizer.SYMBOL {
		e.addError(errors.Errorf("eatSymbol() was given token type %q", e.tz.TokenType()))
		return false
	}
	if e.tz.Symbol() != value {
		e.addError(errors.Errorf("eatSymbol() expect %q, but got %q", value, e.tz.Symbol()))
		return false
	}
	e.advance()
	return true
}

// カレントトークンがSymbolである場合はその値とtrueをかえして次のトークンに進む
// そうでない場合は何もせずエラーを追加してfalseを返す
func (e *Engine) expectSymbol() (string, bool) {
	if e.tz.TokenType() != tokenizer.SYMBOL {
		e.addError(errors.Errorf("expectSymbol() was given token type %q", e.tz.TokenType()))
		return "", false
	}
	value := e.tz.Symbol()
	e.advance()
	return value, true
}

// カレントトークンがOperatorであるばあいはそれを返してtokenをすすめ、trueを返す
// そうでない場合はfalseを返す
func (e *Engine) maybeOp() (string, bool) {
	if e.tz.TokenType() != tokenizer.SYMBOL {
		return "", false
	}
	switch s := e.tz.Symbol(); s {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		e.advance()
		return s, true
	default:
		return "", false
	}
}

// カレントトークンがIdentiferである場合はその値とtrueをかえして次のトークンに進む
// そうでない場合は何もせずエラーを追加してfalseを返す
func (e *Engine) expectIdentifier() (string, bool) {
	if e.tz.TokenType() != tokenizer.IDENTIFIER {
		e.addError(errors.Errorf("expectIdentifier() was given token type %q", e.tz.TokenType()))
		return "", false
	}
	value := e.tz.Identifier()
	e.advance()
	return value, true
}

func (e *Engine) expectIntegerConstant() (int, bool) {
	if e.tz.TokenType() != tokenizer.INT_CONST {
		e.addError(errors.Errorf("expectIdentifier() was given token type %q", e.tz.TokenType()))
		return 0, false
	}
	value := e.tz.IntVal()
	e.advance()
	return value, true
}

func (e *Engine) addError(err error) {
	e.errs = multierror.Append(e.errs, err)
}

func (e *Engine) Error() error {
	return e.errs
}

func (e *Engine) advance() {
	e.tz.Advance()
	// e.logCurrentToken()
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
