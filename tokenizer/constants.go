package tokenizer

type KeyWord string

const (
	CLASS       KeyWord = "class"
	METHOD      KeyWord = "method"
	FUNCTION    KeyWord = "function"
	CONSTRUCTOR KeyWord = "constructor"
	INT         KeyWord = "int"
	BOOLEAN     KeyWord = "boolean"
	CHAR        KeyWord = "char"
	VOID        KeyWord = "void"
	VAR         KeyWord = "var"
	STATIC      KeyWord = "static"
	FIELD       KeyWord = "field"
	LET         KeyWord = "let"
	DO          KeyWord = "do"
	IF          KeyWord = "if"
	ELSE        KeyWord = "else"
	WHILE       KeyWord = "while"
	RETURN      KeyWord = "return"
	TRUE        KeyWord = "true"
	FALSE       KeyWord = "false"
	NULL        KeyWord = "null"
	THIS        KeyWord = "this"
)

func isKeyword(w string) (KeyWord, bool) {
	switch k := KeyWord(w); k {
	case CLASS, METHOD, FUNCTION, CONSTRUCTOR, INT, BOOLEAN, CHAR, VOID, VAR, STATIC, FIELD, LET, DO, IF, ELSE, WHILE, RETURN, TRUE, FALSE, NULL, THIS:
		return k, true
	default:
		return "", false
	}
}
