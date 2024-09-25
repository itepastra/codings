package lzss

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
		israw, err := r.ReadBits(1)
		if err != nil {
			log.Criticalf("error reading rawbit: %e", err)
		}
		if israw == 1 {
			nextchar, err := r.ReadByte()
			if err != nil {
				log.Criticalf("error reading nextchar: %e", err)
				return textBuf
			}
			textBuf = append(textBuf, nextchar)
		} else {
			uoffset, err := r.ReadBits(offsetExp)
			if err != nil {
				log.Criticalf("error reading offset: %e", err)
			}
			offset := int(uoffset)
			ulength, err := r.ReadBits(lengthExp)
			if err != nil {
				log.Criticalf("error reading length: %e", err)
			}
			length := int(ulength)
			nextchar, err := r.ReadByte()
			if err != nil {
				log.Criticalf("error reading nextchar: %e", err)
				return textBuf
			}
			log.Debugf("offset: %d, length %d, next %+q", offset, length, string(nextchar))
			if length != 0 {
				textBuf = append(textBuf, textBuf[(position-offset):(position-offset+length)]...)
				position += length
			}
			textBuf = append(textBuf, nextchar)
			position += 1
			log.Debug(textBuf)
		}
	}
}
