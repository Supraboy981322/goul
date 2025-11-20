package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	input []string
	defsGlob gomn.Map
)

func init() {
	file, err := os.ReadFile("defs.gomn")
	if err != nil {
		log.Fatalf("failed to read defs.gomn:\n  %v", err)
	}
	
	defsGlob, err = gomn.Parse(string(file))
	if err != nil {
		log.Fatal(err)
	}

	input = append(input, "fn prim() {\n  wr.l\n}")
}

func main() {
	input := strings.Split(input[0], " ")
	fmt.Printf("original: %#v\n", input)
	var output []string
	output = parse(input, output, defsGlob, false)

/*	for _, chunk := range output {
		fmt.Print(chunk + " ")
	}*/	fmt.Printf("new: %#v\n", output)
}

func parse(in []string, out []string, defs gomn.Map, sub bool) []string {
	for _, chunk := range in {
		if subFunc := strings.Split(chunk, "."); len(subFunc) > 1 {
			if subDefs, ok := defs[subFunc[0]].(gomn.Map); ok {
				/*if subDefs[0]*/
				out = parse(subFunc, out, subDefs, true)
      } else {
				out = append(out, chunk)
			}
		} else { 
			if newChunk, ok := defs[chunk].(string); ok || sub {
				if sub {
					subChunk, _ := defs[""].(string)
					out = append(out, subChunk)
				} else {
					out = append(out, newChunk)
				}
			} else if sub {
					subChunk, _ := defs[""].(string)
					out = append(out, subChunk)
			} else {
				out = append(out, chunk)
			}
		}
	}
	return out
}
