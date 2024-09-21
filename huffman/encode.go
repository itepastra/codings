package huffman

import (
	"container/heap"
	"fmt"
	"io"
	"log"

	"github.com/bits-and-blooms/bitset"
)

func Encode(text string, writer io.Writer) {
	cm := countMap(text)
	log.Print(cm)

	tree, _ := MkTree(cm)
	log.Print(tree)
	encs := tree.encodings("")
	for _, enc := range encs {
		log.Print(enc)
	}

	mp := tree.makeMap(bitset.New(0))

	for k, v := range mp {
		log.Print(k, v)
	}

	for _, byte := range []byte(text) {
		log.Printf("%s = %s, len %d", string(byte), mp[byte], mp[byte].Len())
		_, _ = (mp[byte]).WriteTo(writer)
	}

}

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

func (h hufTree) encodings(prefix string) []string {
	if h.l != nil && h.r != nil {
		return append(h.l.encodings(fmt.Sprintf("%s0", prefix)), h.r.encodings(fmt.Sprintf("%s1", prefix))...)
	} else {
		return []string{fmt.Sprintf("%s is '%s'", prefix, string(h.val))}
	}
}

func (h hufTree) makeMap(prefix *bitset.BitSet) map[byte]*bitset.BitSet {
	merr := make(map[byte]*bitset.BitSet)

	if h.l != nil && h.r != nil {
		lprefix := prefix.Clone().SetTo(prefix.Len(), true)
		for k, v := range h.l.makeMap(lprefix) {
			merr[k] = v
		}
		rprefix := prefix.Clone().SetTo(prefix.Len(), false)
		for k, v := range h.r.makeMap(rprefix) {
			merr[k] = v
		}
	} else {
		merr[h.val] = prefix
	}

	return merr
}

func MkTree(counts map[byte]int) (hufTree, bool) {

	leafs := make(map[hufTree]int, len(counts))
	for char, amount := range counts {
		leafs[hufTree{char, nil, nil}] = amount
	}

	pq := makeFromMap(leafs)
	log.Print(pq)

	for (*pq).Len() > 1 {
		p1 := heap.Pop(pq).(*Item[hufTree])
		p2 := heap.Pop(pq).(*Item[hufTree])

		newp := hufTree{val: 0, l: &p1.value, r: &p2.value}
		newi := Item[hufTree]{value: newp, priority: p1.priority + p2.priority}

		heap.Push(pq, &newi)
	}

	tree, ok := heap.Pop(pq).(*Item[hufTree])

	return tree.value, ok

}

func countMap(text string) map[byte]int {
	bytemap := make(map[byte]int)
	for i, byte := range []byte(text) {
		log.Printf("byte %d is %s (%b)", i, string(byte), byte)
		bytemap[byte] += 1
	}

	return bytemap
}
