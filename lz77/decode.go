package lz77

import (
	"io"

	"github.com/icza/bitio"
)

func Decode(r *bitio.Reader) []byte {
	offsetExpTemp, err := r.ReadBits(6)
	if err != nil {
		log.Criticalf("offsetExp error %e", err)
	}
	offsetExp := byte(offsetExpTemp)
	lengthExpTemp, err := r.ReadBits(6)
	if err != nil {
		log.Criticalf("lengthExp error %e", err)
	}
	lengthExp := byte(lengthExpTemp)
	log.Debugf("offset bits: %d, length bits: %d", offsetExp, lengthExp)

	textBuf := []byte{}
	position := 0

	for {
		uoffset, err := r.ReadBits(offsetExp)
		if err != nil {
			return textBuf
		}
		offset := int(uoffset)
		ulength, err := r.ReadBits(lengthExp)
		if err != nil {
			return textBuf
		}
		length := int(ulength)
		nextchar, err := r.ReadByte()
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
		} else {
			for i := range length {
				buf = append(buf, buf[position-offset+i])
			}
		}
	}
	return buf
}
