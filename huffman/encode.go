package huffman

import (
	"container/heap"
	"fmt"

	logging "github.com/op/go-logging"

	"github.com/icza/bitio"
)

var log = logging.MustGetLogger("encoder")

func Encode(text string, writer *bitio.Writer) {
	cm := countMap(text)
	log.Infof("count map is: %v", cm)

	tree := MkTree(cm)
	log.Infof("tree is: %v", tree)

	if log.IsEnabledFor(logging.DEBUG) {
		encs := tree.encodings("")
		for _, enc := range encs {
			log.Debug(enc)
		}
	}

	mp := tree.makeMap([]bool{})

	if log.IsEnabledFor(logging.DEBUG) {
		for k, v := range mp {
			log.Debug(k, v)
		}
		for _, byte := range []byte(text) {
			log.Debugf("%s = %v", string(byte), mp[byte])
		}
	}

	// 0 -> level deeper
	// 1 -> value -> byte
	writer.Align()
	tree.treeEncode(writer, true)

	for _, char := range []byte(text) {
		for _, bit := range mp[char] {
			writer.WriteBool(bit)
		}
	}

	_, err := writer.Align()
	if err != nil {
		log.Critical(err)
	}
}

func (h hufTree) treeEncode(w *bitio.Writer, is_first bool) {
	if is_first {
		h.l.treeEncode(w, false)
		h.r.treeEncode(w, false)
	} else if h.l == nil && h.r == nil {
		w.WriteBool(true)
		w.WriteByte(h.val)
	} else {
		w.WriteBool(false)
		h.l.treeEncode(w, false)
		h.r.treeEncode(w, false)
	}
}

func (h hufTree) encodings(prefix string) []string {
	if h.l != nil && h.r != nil {
		return append(h.l.encodings(fmt.Sprintf("%s0", prefix)), h.r.encodings(fmt.Sprintf("%s1", prefix))...)
	} else {
		return []string{fmt.Sprintf("%s is '%s'", prefix, string(h.val))}
	}
}

func (h hufTree) makeMap(prefix []bool) map[byte][]bool {
	merr := make(map[byte][]bool)
	log.Debugf("am at %v", h)

	if h.l == nil && h.r == nil {
		if len(prefix) == 0 {
			log.Debug("tree only has 1 value")
			merr[h.val] = []bool{true}
			return merr
		}
		log.Infof("letter '%s' is %v", string(h.val), prefix)
		merr[h.val] = append([]bool{}, prefix...)
	} else {
		log.Debug("it's a tree")
		lprefix := append(prefix, true)
		for k, v := range h.l.makeMap(lprefix) {
			merr[k] = v
		}
		log.Debug("did the left")
		rprefix := append(prefix, false)
		for k, v := range h.r.makeMap(rprefix) {
			merr[k] = v
		}
		log.Debug("did the right")
	}

	return merr
}

func MkTree(counts map[byte]int) hufTree {

	leafs := make(map[hufTree]int, len(counts))
	for char, amount := range counts {
		leafs[hufTree{char, nil, nil}] = amount
	}

	pq := makeFromMap(leafs)

	if (*pq).Len() == 1 {
		log.Debug("heap only has 1 value")
		return heap.Pop(pq).(*Item[hufTree]).value
	}

	for (*pq).Len() > 1 {
		p1 := heap.Pop(pq).(*Item[hufTree])
		p2 := heap.Pop(pq).(*Item[hufTree])

		newp := hufTree{val: 0, l: &p1.value, r: &p2.value}
		newi := Item[hufTree]{value: newp, priority: p1.priority + p2.priority}

		heap.Push(pq, &newi)
	}

	tree := heap.Pop(pq).(*Item[hufTree])

	return tree.value

}

func countMap(text string) map[byte]int {
	bytemap := make(map[byte]int)
	for i, byte := range []byte(text) {
		log.Debugf("byte %d is %s (%b)", i, string(byte), byte)
		bytemap[byte] += 1
	}

	return bytemap
}
