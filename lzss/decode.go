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
			return textBuf
		}
		if israw == 1 {
			log.Debugf("is raw")
			nextchar, err := r.ReadByte()
			if err != nil {
				log.Critical(err)
				return textBuf
			}
			log.Debugf("israw: %+q", nextchar)
			textBuf = append(textBuf, nextchar)
			position += 1
		} else {
			log.Debugf("is not raw")
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
			log.Debugf("offset: %d, length %d", offset, length)
			if position-offset+length <= len(textBuf) {
				textBuf = append(textBuf, textBuf[(position-offset):(position-offset+length)]...)
			} else {
				for i := range length {
					textBuf = append(textBuf, textBuf[position-offset+i])
				}
			}
			position += length
			log.Debug(textBuf)
		}
	}
}
