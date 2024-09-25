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

const ENCODE_TYPES = 4

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
	} else {
		var file string
		var chosenType string
		var fileText []byte
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewFilePicker().
					Height(10).
					Value(&file).
					Key("original").
					Title("Pick what file you want to encode"),
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
					OptionsFunc(func() []huh.Option[string] {
						options := []huh.Option[string]{}
						ch := make(chan encoded, ENCODE_TYPES)
						wg := sync.WaitGroup{}

						for i := range ENCODE_TYPES {
							wg.Add(1)
							go func(instance int, channel chan encoded, waitgroup *sync.WaitGroup) {
								defer waitgroup.Done()
								var name string
								buf := bytes.NewBuffer([]byte{})
								bitio := bitio.NewWriter(buf)
								switch instance {
								case 0:
									name = "huffman"
									huffman.Encode(fileText, bitio)
								case 1:
									name = "lz77 (16, 8)"
									lz77.Encode(fileText, 16, 8, bitio)
								case 2:
									name = "lz77 (8, 4)"
									lz77.Encode(fileText, 8, 4, bitio)
								case 3:
									name = "lz77 (4, 4)"
									lz77.Encode(fileText, 4, 4, bitio)
								}
								_, err := bitio.Align()
								if err != nil {
									log.Warningf("encoding error with %s (%e)", name, err)
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

						return options
					}, &fileText),
			),
		)

		err := form.Run()
		if err != nil {
			log.Critical(err)
		}
	}
}
