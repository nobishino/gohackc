// Code generated by "stringer -type=NodeKind"; DO NOT EDIT.

package compilation

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[N_CLASS-0]
	_ = x[N_KEYWORD-1]
	_ = x[N_IDENTIFIER-2]
	_ = x[N_SYMBOL-3]
}

const _NodeKind_name = "N_CLASSN_KEYWORDN_IDENTIFIERN_SYMBOL"

var _NodeKind_index = [...]uint8{0, 7, 16, 28, 36}

func (i NodeKind) String() string {
	if i < 0 || i >= NodeKind(len(_NodeKind_index)-1) {
		return "NodeKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeKind_name[_NodeKind_index[i]:_NodeKind_index[i+1]]
}