package tokenizer

import "io"

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
		}
	}
	return nil
}

func symbolTag(symbol string) string {
	return "<symbol> " + symbol + " </symbol>\n"
}
