package lz77

import (
	"io"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/common"
)

func (lz LZ77) Decode(reader io.Reader) []byte {
	defer common.HandleDecodePanic(log)
	bitreader := bitio.NewReader(reader)
	offsetExpTemp, err := bitreader.ReadBits(6)
	if err != nil {
		log.Criticalf("offsetExp error %e", err)
	}
	lz.OffsetBits = byte(offsetExpTemp)
	lengthExpTemp, err := bitreader.ReadBits(6)
	if err != nil {
		log.Criticalf("lengthExp error %e", err)
	}
	lz.LengthBits = byte(lengthExpTemp)
	log.Debugf("offset bits: %d, length bits: %d", lz.OffsetBits, lz.LengthBits)

	textBuf := []byte{}
	position := 0

	for {
		uoffset, err := bitreader.ReadBits(lz.OffsetBits)
		if err != nil {
			return textBuf
		}
		offset := int(uoffset)
		ulength, err := bitreader.ReadBits(lz.LengthBits)
		if err != nil {
			return textBuf
		}
		length := int(ulength)
		nextchar, err := bitreader.ReadByte()
		if err != nil {
			if err == io.EOF {
				textBuf = lzExtend(textBuf, position, offset, length)
			}
			return textBuf
		}
		log.Debugf("offset: %d, length %d, next %+q", offset, length, string(nextchar))
		textBuf = lzExtend(textBuf, position, offset, length)
		textBuf = append(textBuf, nextchar)
		position += length + 1
	}
}

func lzExtend(buf []byte, position int, offset int, length int) []byte {
	if position-offset >= 0 && length != 0 {
		if position-offset+length <= len(buf) {
			buf = append(buf, buf[(position-offset):(position-offset+length)]...)
		} else if position-offset < len(buf) {
			for i := range length {
				buf = append(buf, buf[position-offset+i])
			}
		}
	}
	return buf
}
