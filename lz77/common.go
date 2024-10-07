package lz77

import "github.com/op/go-logging"

var log = logging.MustGetLogger("encoder")

type LZ77 struct {
	OffsetBits byte
	LengthBits byte
}
