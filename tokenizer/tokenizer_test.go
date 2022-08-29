package tokenizer_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/nobishino/gohackc/tokenizer"
)

func TestTokenizer_SmallSource(t *testing.T) {
	testcases := []struct {
		src  string
		want string
	}{
		{"", "<tokens>\n</tokens>"},
		{"\r \t\n\n\t", "<tokens>\n</tokens>"},
		{"{", "<tokens>\n<symbol> { </symbol>\n</tokens>"},
		{"\r\n { /**hello*/\n\r \t//\t }\n", "<tokens>\n<symbol> { </symbol>\n</tokens>"},
		{"{*", "<tokens>\n<symbol> { </symbol>\n<symbol> * </symbol>\n</tokens>"},
		{"12345\n", "<tokens>\n<integerConstant> 12345 </integerConstant>\n</tokens>"},
		{`"Hello" \n\t"World"\n`, "<tokens>\n<stringConstant> Hello </stringConstant>\n<stringConstant> World </stringConstant>\n</tokens>"},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			tnzr := tokenizer.NewTokenizer(strings.NewReader(tc.src))
			var dst strings.Builder
			if err := tokenizer.ToXML(&dst, tnzr); err != nil {
				t.Fatal(err)
			}
			got := dst.String()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Tokenize result(-want +got):\n%s", diff)
			}
		})
	}
}

// func TestTokenizer_RealSource(t *testing.T) {
// 	testcases := []struct {
// 		srcFile  string
// 		wantFile string
// 	}{
// 		{"ArrayTest/Main.jack", "ArrayTest/MainT.xml"},
// 	}
// 	const testDir = "../testdata"
// 	for _, tc := range testcases {
// 		t.Run(tc.srcFile, func(t *testing.T) {
// 			s := openFile(t, filepath.Join(testDir, tc.srcFile))
// 			tnzr := tokenizer.NewTokenizer(s)
// 			var dst strings.Builder

// 			if err := tokenizer.ToXML(&dst, tnzr); err != nil {
// 				t.Fatal(err)
// 			}
// 			got := dst.String()
// 			want := readFile(t, filepath.Join(testDir, tc.wantFile))
// 			if diff := cmp.Diff(want, got); diff != "" {
// 				t.Errorf("Tokenize result(-want +got):\n%s", diff)
// 			}

// 		})
// 	}
// }

// func openFile(t *testing.T, path string) *os.File {
// 	t.Helper()

// 	f, err := os.Open(path)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Cleanup(func() { f.Close() })
// 	return f
// }

// func readFile(t *testing.T, path string) string {
// 	f := openFile(t, path)
// 	var sb strings.Builder
// 	if _, err := io.Copy(&sb, f); err != nil {
// 		t.Fatal(err)
// 	}
// 	return sb.String()
// }
