package tokenizer

type keyWord string

const (
	CLASS       keyWord = "CLASS"
	METHOD      keyWord = "METHOD"
	FUNCTION    keyWord = "FUNCTION"
	CONSTRUCTOR keyWord = "CONSTRUCTOR"
	INT         keyWord = "INT"
	BOOLEAN     keyWord = "BOOLEAN"
	CHAR        keyWord = "CHAR"
	VOID        keyWord = "VOID"
	VAR         keyWord = "VAR"
	STATIC      keyWord = "STATIC"
	FIELD       keyWord = "FIELD"
	LET         keyWord = "LET"
	DO          keyWord = "DO"
	IF          keyWord = "IF"
	ELSE        keyWord = "ELSE"
	WHILE       keyWord = "WHILE"
	RETURN      keyWord = "RETURN"
	TRUE        keyWord = "TRUE"
	FALSE       keyWord = "FALSE"
	NULL        keyWord = "NULL"
	THIS        keyWord = "THIS"
)
