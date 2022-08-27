package tokenizer

type token struct {
	tokenType
	symbol      string
	keyword     string
	intValue    int
	stringValue string
	identifier  string
}
type tokenType string

const (
	KEYWORD      tokenType = "KEYWORD"
	SYMBOL       tokenType = "SYMBOL"
	IDENTIFIER   tokenType = "IDENTIFIER"
	INT_CONST    tokenType = "INT_CONST"
	STRING_CONST tokenType = "STRING_CONST"
	EOF          tokenType = "EOF"
)
