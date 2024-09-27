package lzss

import (
	"fmt"
	"math"
	"slices"

	"github.com/icza/bitio"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("encoder")

type datatriplet struct {
	israw   bool
	data    byte
	triplet duo
}

type duo struct {
	offset int
	length int
}

func Encode(text []byte, searchBufferLengthExp byte, lookaheadLengthExp byte, writer *bitio.Writer) {
	if searchBufferLengthExp > 63 {
		log.Criticalf("searchbuffer exponent (%d) not allowed", searchBufferLengthExp)
		return
	}
	if lookaheadLengthExp > 63 {
		log.Criticalf("lookahead exponent (%d) not allowed", lookaheadLengthExp)
		return
	}
	textBytes := []byte(text)
	searchStart := 0
	searchEnd := 0
	codingPosition := 0
	encodedTriplets := []datatriplet{}

	searchBufferLength := 1 << int(searchBufferLengthExp)
	lookaheadLength := 1 << int(lookaheadLengthExp)
	breakeven := int(math.Ceil((float64(searchBufferLengthExp) + float64(lookaheadLengthExp)) / 8.0))

	for codingPosition < len(textBytes) {
		lookahead_end := min(codingPosition+lookaheadLength+1, len(textBytes))
		lookahead := textBytes[codingPosition:lookahead_end]
		offset, length, nextChar := findMatch(textBytes[searchStart:searchEnd], lookahead, lookaheadLength)

		if nextChar != nil {
			log.Debugf("nextchar is %+q", string(*nextChar))
			var newTriplet datatriplet
			if length > breakeven {
				newTriplet = datatriplet{
					israw: false,
					triplet: duo{
						offset: offset,
						length: length,
					},
				}
				codingPosition += length
				searchEnd += length
			} else {
				newTriplet = datatriplet{
					israw: true,
					data:  lookahead[0],
				}
				codingPosition += 1
				searchEnd += 1
			}
			encodedTriplets = append(encodedTriplets, newTriplet)
		} else {
			var newTriplet datatriplet
			if length > breakeven {
				newTriplet = datatriplet{
					israw: false,
					triplet: duo{
						offset: offset,
						length: length,
					}}
				codingPosition += length
				searchEnd += length
			} else {
				newTriplet = datatriplet{
					israw: true,
					data:  lookahead[0],
				}
				codingPosition += 1
				searchEnd += 1
			}
			encodedTriplets = append(encodedTriplets, newTriplet)
		}
		searchStart = max(0, searchEnd-searchBufferLength)
	}

	log.Infof("triplets %v", encodedTriplets)

	writer.WriteBits(uint64(searchBufferLengthExp), 6)
	writer.WriteBits(uint64(lookaheadLengthExp), 6)
	for _, dtri := range encodedTriplets {
		if dtri.israw {
			err := writer.WriteBits(1, 1)
			if err != nil {
				log.Critical(err)
			}
			err = writer.WriteByte(dtri.data)
			if err != nil {
				log.Critical(err)
			}
		} else {
			tri := dtri.triplet
			log.Debugf("writing offset %d, writing length %d", tri.offset, tri.length)
			err := writer.WriteBits(0, 1)
			if err != nil {
				log.Critical(err)
			}
			err = writer.WriteBits(uint64(tri.offset), searchBufferLengthExp)
			if err != nil {
				log.Critical(err)
			}
			err = writer.WriteBits(uint64(tri.length), lookaheadLengthExp)
			if err != nil {
				log.Critical(err)
			}
		}
	}

	_, err := writer.Align()
	if err != nil {
		log.Critical(err)
	}
	return
}

func (t duo) String() string {
	return fmt.Sprintf("[ %d, %d ]", t.offset, t.length)
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
