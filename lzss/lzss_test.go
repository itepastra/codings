package lzss_test

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/lzss"
	testinghelpers "github.com/itepastra/codings/testing_helpers"
)

func TestLzssEncodeDecodeShort(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	for range 100 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encodeBits := bitio.NewWriter(encoded)
	lzss.Encode(buf.Bytes(), 16, 8, encodeBits)

	decodeBits := bitio.NewReader(encoded)
	decoded := lzss.Decode(decodeBits)

	if len(decoded) != len(buf.Bytes()) {
		t.Logf("decoded had length %d, while the buffer had length %d", len(decoded), buf.Len())
		t.FailNow()
	}

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}

func TestLzssEncodeDecodeLong(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	for range 1000 {
		buf.WriteByte(byte(rand.N(93) + 33))
	}
	t.Logf("the input text is: %s", string(buf.Bytes()))

	encoded := bytes.NewBuffer([]byte{})
	encodeBits := bitio.NewWriter(encoded)
	lzss.Encode(buf.Bytes(), 16, 8, encodeBits)

	decodeBits := bitio.NewReader(encoded)
	decoded := lzss.Decode(decodeBits)

	if len(decoded) != len(buf.Bytes()) {
		t.Logf("decoded had length %d, while the buffer had length %d", len(decoded), buf.Len())
		t.FailNow()
	}

	for i, byte := range buf.Bytes() {
		testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
	}
}

var explicits = [][]byte{
	[]byte("aaaaaaaaaaaaaaaaaa"),
}

func TestLzssEncodeDecodeSet(t *testing.T) {
	for _, explicit := range explicits {
		encoded := bytes.NewBuffer([]byte{})
		encodeBits := bitio.NewWriter(encoded)
		lzss.Encode(explicit, 16, 8, encodeBits)

		decodeBits := bitio.NewReader(encoded)
		decoded := lzss.Decode(decodeBits)

		if len(decoded) != len(explicit) {
			t.Logf("decoded had length %d, while the buffer had length %d", len(decoded), len(explicit))
			t.FailNow()
		}

		for i, byte := range explicit {
			testinghelpers.ExpectEqual(t, decoded[i], byte, fmt.Sprint(i))
		}
	}
}
