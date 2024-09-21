package main

import (
	"flag"
	"os"

	"github.com/itepastra/codings/huffman"
)

var t = flag.String("encoding", "huffman", "what encoding to use")

func main() {
	if *t == "huffman" {
		text := "If I write some text here, I wonder what happens to the huffman tree. While I assume it's correct I'm not certain"
		huffman.Encode(text, os.Stdout)
	}
}
