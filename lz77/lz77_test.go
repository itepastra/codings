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

var explicits = [][]byte{
	[]byte("aaaaa"),
}

func TestLz77EncodeDecodeSet(t *testing.T) {
	t.Parallel()
	for _, explicit := range explicits {
		encoded := bytes.NewBuffer([]byte{})
		encoder := lz77.LZ77{OffsetBits: 16, LengthBits: 8}
		encoder.Encode(explicit, encoded)

		decoder := lz77.LZ77{}
		decodeBits := bitio.NewReader(encoded)
		decoded := decoder.Decode(decodeBits)

		if len(decoded) != len(explicit) {
			t.Logf("decoded had length %d, while the buffer had length %d", len(decoded), len(explicit))
			t.FailNow()
		}

		for i, byte := range explicit {
			testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
		}
	}
}

func TestLz77EncodeDecodeShort(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer([]byte{})
	for range 100 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encoder := lz77.LZ77{OffsetBits: 16, LengthBits: 8}
	encoder.Encode(buf.Bytes(), encoded)

	decoder := lz77.LZ77{}
	decodeBits := bitio.NewReader(encoded)
	decoded := decoder.Decode(decodeBits)

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}

func TestLz77EncodeDecodeLong(t *testing.T) {
	t.Parallel()
	buf := bytes.NewBuffer([]byte{})
	for range 10000 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encoder := lz77.LZ77{OffsetBits: 16, LengthBits: 8}
	encoder.Encode(buf.Bytes(), encoded)

	decoder := lz77.LZ77{}
	decodeBits := bitio.NewReader(encoded)
	decoded := decoder.Decode(decodeBits)

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}
