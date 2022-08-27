package tokenizer

import "io"

func ToXML(dst io.Writer, t *Tokenizer) error {
loop:
	for t.HasMoreTokens() {
		t.Advance()
		switch t.currentToken.tokenType {
		case EOF:
			break loop
		}
	}
	return nil
}
