package lz77

import (
	"fmt"
	"io"
	"slices"

	"github.com/icza/bitio"
)

type triplet struct {
	offset   int
	length   int
	next     bool
	nextChar byte
}

func (lz LZ77) Encode(text []byte, writer io.Writer) {
	if lz.OffsetBits > 63 {
		log.Criticalf("searchbuffer exponent (%d) not allowed", lz.OffsetBits)
		return
	}
	if lz.LengthBits > 63 {
		log.Criticalf("lookahead exponent (%d) not allowed", lz.LengthBits)
		return
	}

	bitwriter := bitio.NewWriter(writer)
	defer bitwriter.Close()

	textBytes := []byte(text)
	searchStart := 0
	searchEnd := 0
	codingPosition := 0
	encodedTriplets := []triplet{}

	searchBufferLength := 1 << int(lz.OffsetBits)
	lookaheadLength := 1 << int(lz.LengthBits)

	for codingPosition < len(textBytes) {
		lookahead_end := min(codingPosition+lookaheadLength+1, len(textBytes))
		lookahead := textBytes[codingPosition:lookahead_end]
		offset, length, nextChar := findMatch(textBytes[searchStart:searchEnd], lookahead, lookaheadLength)

		codingPosition += length + 1
		searchEnd += length + 1
		searchStart = max(0, searchEnd-searchBufferLength)
		if codingPosition > len(textBytes) {
			encodedTriplets = append(encodedTriplets, triplet{offset: offset, length: length, next: false})
		} else if nextChar != nil {
			log.Debugf("nextchar is %+q", string(*nextChar))
			encodedTriplets = append(encodedTriplets, triplet{offset: offset, length: length, next: true, nextChar: *nextChar})
		} else {
			log.Critical("somehow there was not nextcharacter but we're not at the end")
			encodedTriplets = append(encodedTriplets, triplet{offset: offset, length: length, next: false})
		}
	}

	log.Infof("triplets %v", encodedTriplets)

	bitwriter.WriteBits(uint64(lz.OffsetBits), 6)
	bitwriter.WriteBits(uint64(lz.LengthBits), 6)
	for _, tri := range encodedTriplets {
		err := bitwriter.WriteBits(uint64(tri.offset), lz.OffsetBits)
		if err != nil {
			log.Critical(err)
		}
		err = bitwriter.WriteBits(uint64(tri.length), lz.LengthBits)
		if err != nil {
			log.Critical(err)
		}
		if tri.next {
			err := bitwriter.WriteByte(tri.nextChar)
			if err != nil {
				log.Critical(err)
			}
		} else {
			log.Debug("returning because there is no next")
			_, err := bitwriter.Align()
			if err != nil {
				log.Critical(err)
			}
			return
		}
	}

	log.Debug("returning because done")
	_, err := bitwriter.Align()
	if err != nil {
		log.Critical(err)
	}
	return
}

func (t triplet) String() string {
	if t.next {
		return fmt.Sprintf("[ %d, %d, %+q ]", t.offset, t.length, string(t.nextChar))
	} else {
		return fmt.Sprintf("[ %d, %d, - ]", t.offset, t.length)
	}
}

func findMatch(search []byte, lookahead []byte, lookaheadLength int) (offset int, length int, next *byte) {

	if !slices.Contains(search, lookahead[0]) {
		return 0, 0, &lookahead[0]
	}

	n := 0
	for i := len(search) - 1; i >= 0; i -= 1 {
		n += 1
		for j := 0; cmp(search, lookahead, i, j); j += 1 {
			m := j + 1
			if m > length {
				length = m
				offset = n
				if len(lookahead) > m {
					next = &lookahead[m]
					log.Debugf("new longest match: %v %+q, %d, %d", next, string(*next), length, offset)
					if lookaheadLength == m {
						return
					}
				}
			}
		}
	}

	log.Debugf("next addr is %v", next)
	return
}

func cmp(search []byte, lookahead []byte, index int, offset int) bool {
	searchlen := len(search)
	var val byte
	if index+offset >= searchlen {
		val = lookahead[index+offset-searchlen]
	} else {
		val = search[index+offset]
	}
	if offset < len(lookahead) {
		return val == lookahead[offset]
	} else {
		return false
	}
}
