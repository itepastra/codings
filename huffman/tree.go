package huffman

import "fmt"

type hufTree struct {
	val byte
	l   *hufTree
	r   *hufTree
}

func (h hufTree) String() string {
	if h.l != nil && h.r != nil {
		return fmt.Sprintf("(%s) (%s)", h.l, h.r)
	} else {
		return fmt.Sprintf("%s", string(h.val))
	}
}
