// Package compilation はコンパイルを行う
package compilation

import "io"

type Engine struct {
	src io.Reader
	dst io.Writer
}

func New(src io.Reader, dst io.Writer) *Engine {
	return &Engine{
		src: src,
		dst: dst,
	}
}
