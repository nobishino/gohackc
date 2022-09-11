package analyzer

import (
	"os"

	"github.com/pkg/errors"
)

// Analyze はXxx.jackからXxx.xmlをつくる
// source: ファイル名かディレクトリ名
// 一旦ファイルにのみ対応
func Analyze(source string) error {
	s, err := os.Open(source)
	if err != nil {
		return errors.WithStack(err)
	}
	defer s.Close()
	// e := compilation.New(s)
	// tree, err := e.CompileClass()
	// if err != nil {
	// 	return err
	// }
	return nil
}
