package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/icza/bitio"
	"github.com/itepastra/codings/common"
	"github.com/itepastra/codings/huffman"
	"github.com/itepastra/codings/lz77"
	"github.com/itepastra/codings/lzss"
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

const ENCODE_TYPES = 7
const DECODE_TYPES = 3

type encoded struct {
	codec string
	data  []byte
}

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
		var encoder common.Encoder
		switch *codec {
		case "huffman":
			encoder = huffman.Huffman{}
			encoder.Encode(text, writer)
		case "lz77":
			encoder = lz77.LZ77{OffsetBits: byte(*offsetBits), LengthBits: byte(*lengthBits)}
			encoder.Encode(text, writer)
		case "lzss":
			bitWriter := bitio.NewWriter(writer)
			lzss.Encode(text, byte(*offsetBits), byte(*lengthBits), bitWriter)
			err = bitWriter.Close()
			if err != nil {
				log.Critical(err)
			}
		default:
			log.Warning("encoding not supported")
			return
		}

		writer.Flush()
	} else if *decode {
		reader := bufio.NewReader(os.Stdin)
		writer := bufio.NewWriter(os.Stdout)
		var decoder common.Decoder
		switch *codec {
		case "huffman":
			decoder = huffman.Huffman{}
			writer.Write(decoder.Decode(reader))
		case "lz77":
			decoder = lz77.LZ77{}
			writer.Write(decoder.Decode(reader))
		case "lzss":
			bitreader := bitio.NewReader(reader)
			writer.Write(lzss.Decode(bitreader))
		default:
			log.Warningf("%s encoding not (yet) supported", *codec)
			return
		}
		writer.Flush()
	} else {
		var compress bool
		var file string
		var chosenType string
		var fileText []byte
		var outputPath string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[bool]().
					Height(4).
					Title("What operation?").
					Value(&compress).
					Options(
						huh.NewOption("encode", true),
						huh.NewOption("decode", false),
					),
				huh.NewFilePicker().
					Height(10).
					Value(&file).
					Key("original").
					Title("what file is your input"),
				huh.NewNote().TitleFunc(func() string {
					f, err := os.Open(file)
					if err != nil {
						return "There was an error opening the file"
					}
					defer f.Close()

					text, err := io.ReadAll(f)
					if err != nil {
						return "there was an error reading the file"
					}

					fileText = text
					return fmt.Sprintf("%s has size %dB", file, len(fileText))
				}, &file),
				huh.NewSelect[string]().
					Value(&chosenType).
					Height(10).
					OptionsFunc(func() []huh.Option[string] { return compressionOptions(fileText, compress) }, &fileText),
				huh.NewInput().
					Title("Where to save the output").Value(&outputPath),
			),
		)

		err := form.Run()
		if err != nil {
			log.Critical(err)
			return
		}

		err = os.WriteFile(outputPath, []byte(chosenType), 0644)
		if err != nil {
			log.Critical(err)
			return
		}
	}
}

func compressionOptions(fileText []byte, compress bool) []huh.Option[string] {
	options := []huh.Option[string]{}
	if compress {
		ch := make(chan encoded, ENCODE_TYPES)
		wg := sync.WaitGroup{}

		for i := range ENCODE_TYPES {
			wg.Add(1)
			go func(instance int, channel chan encoded, waitgroup *sync.WaitGroup) {
				defer waitgroup.Done()
				var name string
				var encoder common.Encoder
				buf := bytes.NewBuffer([]byte{})
				switch instance {
				case 0:
					name = "huffman"
					encoder = huffman.Huffman{}
					encoder.Encode(fileText, buf)
				case 1:
					name = "lz77 (16, 8)"
					encoder = lz77.LZ77{OffsetBits: 16, LengthBits: 8}
					encoder.Encode(fileText, buf)
				case 2:
					name = "lz77 (8, 4)"
					encoder = lz77.LZ77{OffsetBits: 8, LengthBits: 4}
					encoder.Encode(fileText, buf)
				case 3:
					name = "lz77 (4, 4)"
					encoder = lz77.LZ77{OffsetBits: 4, LengthBits: 4}
					encoder.Encode(fileText, buf)
				case 4:
					name = "lzss (16, 8)"
					bitio := bitio.NewWriter(buf)
					lzss.Encode(fileText, 16, 8, bitio)
					_, err := bitio.Align()
					if err != nil {
						log.Warningf("encoding error with %s (%e)", name, err)
					}
				case 5:
					name = "lzss (8, 4)"
					bitio := bitio.NewWriter(buf)
					lzss.Encode(fileText, 8, 4, bitio)
					_, err := bitio.Align()
					if err != nil {
						log.Warningf("encoding error with %s (%e)", name, err)
					}
				case 6:
					name = "lzss (4, 4)"
					bitio := bitio.NewWriter(buf)
					lzss.Encode(fileText, 4, 4, bitio)
					_, err := bitio.Align()
					if err != nil {
						log.Warningf("encoding error with %s (%e)", name, err)
					}
				}
				channel <- encoded{name, buf.Bytes()}
				return
			}(i, ch, &wg)
		}

		wg.Wait()
		close(ch)
		for enc := range ch {
			options = append(options, huh.NewOption(fmt.Sprintf("%s (%dB)", enc.codec, len(enc.data)), string(enc.data)))
		}

	} else {

		ch := make(chan encoded, DECODE_TYPES)
		wg := sync.WaitGroup{}

		for i := range DECODE_TYPES {
			wg.Add(1)
			go func(instance int, channel chan encoded, waitgroup *sync.WaitGroup) {
				defer waitgroup.Done()
				var name string
				var output []byte
				var decoder common.Decoder
				buf := bytes.NewBuffer(fileText)
				bitio := bitio.NewReader(buf)
				switch instance {
				case 0:
					name = "huffman"
					decoder = huffman.Huffman{}
					output = decoder.Decode(buf)
				case 1:
					name = "lz77"
					decoder = lz77.LZ77{}
					output = decoder.Decode(buf)
				case 2:
					name = "lzss"
					output = lzss.Decode(bitio)
				}
				channel <- encoded{name, output}
				return
			}(i, ch, &wg)
		}

		wg.Wait()
		close(ch)

		for enc := range ch {
			options = append(options, huh.NewOption(fmt.Sprintf("%s (%dB)", enc.codec, len(enc.data)), string(enc.data)))
		}

	}
	return options
}
