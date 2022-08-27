package tokenizer

type TokenType string

const (
	KEYWORD      TokenType = "KEYWORD"
	SYMBOL       TokenType = "SYMBOL"
	IDENTIFIER   TokenType = "IDENTIFIER"
	INT_CONST    TokenType = "INT_CONST"
	STRING_CONST TokenType = "STRING_CONST"
)

type Keyword string

const (
	CLASS       Keyword = "CLASS"
	METHOD      Keyword = "METHOD"
	FUNCTION    Keyword = "FUNCTION"
	CONSTRUCTOR Keyword = "CONSTRUCTOR"
	INT         Keyword = "INT"
	BOOLEAN     Keyword = "BOOLEAN"
	CHAR        Keyword = "CHAR"
	VOID        Keyword = "VOID"
	VAR         Keyword = "VAR"
	STATIC      Keyword = "STATIC"
	FIELD       Keyword = "FIELD"
	LET         Keyword = "LET"
	DO          Keyword = "DO"
	IF          Keyword = "IF"
	ELSE        Keyword = "ELSE"
	WHILE       Keyword = "WHILE"
	RETURN      Keyword = "RETURN"
	TRUE        Keyword = "TRUE"
	FALSE       Keyword = "FALSE"
	NULL        Keyword = "NULL"
	THIS        Keyword = "THIS"
)
