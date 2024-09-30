package common

import "io"

type Decoder interface {
	Decode(reader io.Reader) []byte
}
