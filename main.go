package main

import (
	"bufio"
	"flag"
	"io"
	"os"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/huffman"
	"github.com/itepastra/codings/lz77"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("encoder")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var encode = flag.Bool("encode", false, "should encode input text")
var decode = flag.Bool("decode", false, "should decode input text")
var codec = flag.String("encoding", "huffman", "what encoding to use")
var lengthBits = flag.Int("length", 8, "how many bits should be used for the length in lz77 (up to 63)")
var offsetBits = flag.Int("offset", 16, "how many bits should be used for the offset in lz77 (up to 63)")

var debug = flag.Bool("debug", false, "should enable debug output")

func main() {
	flag.Parse()
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	if *debug {
		logging.SetLevel(logging.DEBUG, "encoder")
	} else {
		logging.SetLevel(logging.WARNING, "encoder")
	}

	if *encode {
		text, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Warning(err)
		}
		writer := bufio.NewWriter(os.Stdout)
		bitWriter := bitio.NewWriter(writer)
		switch *codec {
		case "huffman":
			huffman.Encode(text, bitWriter)
		case "lz77":
			lz77.Encode(text, byte(*offsetBits), byte(*lengthBits), bitWriter)
		default:
			log.Warning("encoding not supported")
			return
		}

		err = bitWriter.Close()
		if err != nil {
			log.Critical(err)
		}
		writer.Flush()
	} else if *decode {
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout)
		bitreader := bitio.NewReader(reader)
		switch *codec {
		case "huffman":
			writer.Write(huffman.Decode(bitreader))
		case "lz77":
			writer.Write(lz77.Decode(bitreader))
		default:
			log.Warningf("%s encoding not (yet) supported", *codec)
			return
		}
		writer.Flush()
	}
}
