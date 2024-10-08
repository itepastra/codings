package huffman_test

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/huffman"
	testinghelpers "github.com/itepastra/codings/testing_helpers"
)

func TestHufEncodeDecode(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer([]byte{})
	for range 10000 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}

	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoder := huffman.Huffman{}
	encoded := bytes.NewBuffer([]byte{})
	encodeBits := bitio.NewWriter(encoded)
	encoder.Encode(buf.Bytes(), encodeBits)

	decoder := huffman.Huffman{}
	decodeBits := bitio.NewReader(encoded)
	decoded := decoder.Decode(decodeBits)

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}
