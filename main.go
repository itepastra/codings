package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/icza/bitio"
	"github.com/itepastra/codings/huffman"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("encoder")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var encode = flag.Bool("encode", false, "should encode input text")
var decode = flag.Bool("decode", false, "should decode input text")
var t = flag.String("encoding", "huffman", "what encoding to use")

func main() {
	flag.Parse()
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	if *encode {
		bitWriter := bitio.NewWriter(writer)
		text, err := reader.ReadString(0x00)
		if err != nil {
			log.Warning(err)
		}
		if text[len(text)-1] == '\n' {
			text = text[:len(text)-1]
		}
		switch *t {
		case "huffman":
			huffman.Encode(text, bitWriter)
		default:
			log.Warning("encoding not supported")
			return
		}

		err = bitWriter.Close()
		if err != nil {
			log.Critical(err)
		}
		writer.Flush()
	}

	if *decode {
		bitreader := bitio.NewReader(reader)
		switch *t {
		case "huffman":
			writer.Write(huffman.Decode(bitreader))
		default:
			log.Warning("encoding not supported")
			return
		}
		writer.Flush()
	}
}
