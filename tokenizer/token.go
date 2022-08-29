package tokenizer

import "strconv"

type token struct {
	tokenType
	symbol      string
	keyword     keyWord
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

func (tk token) String() string {
	result := string(tk.tokenType)
	var value string
	switch tk.tokenType {
	case EOF:
		value = "EOF"
	case KEYWORD:
		value = string(tk.keyword)
	case SYMBOL:
		value = tk.symbol
	case IDENTIFIER:
		value = tk.identifier
	case INT_CONST:
		value = strconv.Itoa(tk.intValue)
	case STRING_CONST:
		value = tk.stringValue
	}
	return result + ": " + value
}
