package lz77

import "github.com/icza/bitio"

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
			return textBuf
		}
		log.Debugf("offset: %d, length %d, next %+q", offset, length, string(nextchar))
		if length != 0 && position-offset >= 0 && position-offset+length < len(textBuf) {
			textBuf = append(textBuf, textBuf[(position-offset):(position-offset+length)]...)
			position += length
		}
		textBuf = append(textBuf, nextchar)
		position += 1
		log.Debug(textBuf)
	}
}
