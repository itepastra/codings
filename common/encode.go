package common

import "io"

type Encoder interface {
	Encode(text []byte, writer io.Writer)
}
