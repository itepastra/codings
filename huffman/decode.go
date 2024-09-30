package huffman

import (
	"io"

	"github.com/icza/bitio"
)

func (h Huffman) Decode(r io.Reader) []byte {
	bitreader := bitio.NewReader(r)
	tree, err := gentree(bitreader)
	if err != nil {
		log.Critical("decoding error")
	}

	text := tree.decode(bitreader)
	log.Info(text)

	return text
}

func gentree(r *bitio.Reader) (*hufTree, error) {
	selector, err := r.ReadBits(1)
	if err != nil {
		return nil, err
	}
	tree := hufTree{}
	// get left value
	if selector == 1 {
		log.Debug("is value on left")
		nextbyte, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		tree.l = &hufTree{val: nextbyte}
	} else {
		log.Debug("is tree on left")
		t, err := gentree(r)
		if err != nil {
			return nil, err
		}
		tree.l = t
	}
	selector, err = r.ReadBits(1)
	if err != nil {
		return nil, err
	}
	// get right value
	if selector == 1 {
		log.Debug("is value on right")
		nextbyte, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		tree.r = &hufTree{val: nextbyte}
	} else {
		log.Debug("is tree on right")
		t, err := gentree(r)
		if err != nil {
			return nil, err
		}
		tree.r = t
	}
	log.Debugf("have tree %v", tree)
	return &tree, nil
}

func (h hufTree) decode(r *bitio.Reader) []byte {
	localtree := &h
	text := []byte{}
	for {
		bit, err := r.ReadBits(1)
		log.Debugf("bit is %d", bit)
		if err != nil {
			if err != io.EOF {
				log.Critical(err)
			}
			return text
		}
		if bit == 1 {
			localtree = localtree.l
			if localtree.l != nil {
				log.Debug("is not end")
			} else {
				text = append(text, localtree.val)
				log.Debugf("is end, is %s", string(localtree.val))
				localtree = &h
			}
		} else {
			localtree = localtree.r
			if localtree.r != nil {
				log.Debug("is not end")
			} else {
				text = append(text, localtree.val)
				log.Debugf("is end, is %s", string(localtree.val))
				localtree = &h
			}
		}
	}
}
