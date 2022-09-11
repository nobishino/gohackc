package compilation

type Node struct {
	Kind  NodeKind
	Value string
	Lhs   *Node
	Rhs   *Node
}
