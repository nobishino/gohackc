package tokenizer

import (
	"fmt"
	"io"
)

func ToXML(dst io.Writer, t *Tokenizer) error {
	if _, err := io.WriteString(dst, "<tokens>\n"); err != nil {
		return err
	}
	defer io.WriteString(dst, "</tokens>")
loop:
	for t.HasMoreTokens() {
		t.Advance()
		switch t.TokenType() {
		case EOF:
			break loop
		case SYMBOL:
			if _, err := io.WriteString(dst, symbolTag(t.Symbol())); err != nil {
				return err
			}
		case INT_CONST:
			if _, err := io.WriteString(dst, integerConstantTag(t.IntVal())); err != nil {
				return err
			}
		}
	}
	return nil
}

func symbolTag(symbol string) string {
	return "<symbol> " + symbol + " </symbol>\n"
}

func integerConstantTag(value int) string {
	return fmt.Sprintf("<integerConstant> %d </integerConstant>\n", value)
}
