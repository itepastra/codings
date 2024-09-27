package lz77_test

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/lz77"
	testinghelpers "github.com/itepastra/codings/testing_helpers"
)

func TestLz77EncodeDecodeShort(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	for range 100 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encodeBits := bitio.NewWriter(encoded)
	lz77.Encode(buf.Bytes(), 16, 8, encodeBits)

	decodeBits := bitio.NewReader(encoded)
	decoded := lz77.Decode(decodeBits)

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}

func TestLz77EncodeDecodeLong(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	for range 10000 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encodeBits := bitio.NewWriter(encoded)
	lz77.Encode(buf.Bytes(), 16, 8, encodeBits)

	decodeBits := bitio.NewReader(encoded)
	decoded := lz77.Decode(decodeBits)

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}
