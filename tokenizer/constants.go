package tokenizer

type keyWord string

const (
	CLASS       keyWord = "class"
	METHOD      keyWord = "method"
	FUNCTION    keyWord = "function"
	CONSTRUCTOR keyWord = "constructor"
	INT         keyWord = "int"
	BOOLEAN     keyWord = "boolean"
	CHAR        keyWord = "char"
	VOID        keyWord = "void"
	VAR         keyWord = "var"
	STATIC      keyWord = "static"
	FIELD       keyWord = "field"
	LET         keyWord = "let"
	DO          keyWord = "do"
	IF          keyWord = "if"
	ELSE        keyWord = "else"
	WHILE       keyWord = "while"
	RETURN      keyWord = "return"
	TRUE        keyWord = "true"
	FALSE       keyWord = "false"
	NULL        keyWord = "null"
	THIS        keyWord = "this"
)

func isKeyword(w string) (keyWord, bool) {
	switch k := keyWord(w); k {
	case CLASS, METHOD, FUNCTION, CONSTRUCTOR, INT, BOOLEAN, CHAR, VOID, VAR, STATIC, FIELD, LET, DO, IF, ELSE, WHILE, RETURN, TRUE, FALSE, NULL, THIS:
		return k, true
	default:
		return "", false
	}
}
