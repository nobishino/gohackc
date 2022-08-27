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
		switch t.currentToken.tokenType {
		case EOF:
			break loop
		case SYMBOL:
			if _, err := io.WriteString(dst, symbolTag(t.currentToken)); err != nil {
				return err
			}
		}
	}
	return nil
}

func symbolTag(tk token) string {
	return "<symbol> " + tk.symbol + " </symbol>\n"
}
