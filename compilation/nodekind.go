package compilation

type NodeKind int

//go:generate stringer -type=NodeKind

const (
	N_CLASS NodeKind = iota
	N_KEYWORD
	N_IDENTIFIER
	N_SYMBOL
)
