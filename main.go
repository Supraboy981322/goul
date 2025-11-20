package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	isString bool
	killOnWarn bool
	input []string
	inputHeader []string
	splitters []string
	headEnd string
	splitHead bool
	splitHeadAt int
	defsGlob gomn.Map
	rcDefs gomn.Map	
	importsMap gomn.Map
	importDefs gomn.Map
	headDefs gomn.Map
)

func init() {
	var ok bool

	defsFile, err := os.ReadFile("defs.gomn")
	if err != nil {
		log.Fatalf("failed to read defs.gomn:\n  %v", err)
	}
	
	if defsGlob, err = gomn.Parse(string(defsFile)); err != nil {
		log.Fatal(err)
	}

	if rcDefs, ok = defsGlob[0].(gomn.Map); !ok {
		log.Warn("rc definitions not found, may produce odd results")
		kilOcont("continuing anyways")
	} else {
		if killOnWarn, _ = rcDefs["kill on warn"].(bool); killOnWarn {
			log.Info("configured to kill on warn")
		}
		if headEnd, ok = rcDefs["head end"].(string); !ok {
			log.Warn("\"head end\" not defined in rc definitions, this will probably cause problems")
			kilOcont("continuing anyways")
		}
	}
	
	inFile, err := os.ReadFile("foo.gocl")
	if err != nil {
		log.Fatal(err)
	}
	
	strFile := string(inFile)
	var trimmedFile []string
	for _, line := range strings.Split(strFile, "\n") {
		if trimmedLine := strings.TrimSpace(line); trimmedLine != headEnd {
			trimmedFile = append(trimmedFile, trimmedLine)
		} else {
			inputHeader = strings.FieldsFunc(strings.Join(trimmedFile, "\n"), whitespaceSplitter)
			trimmedFile = []string{""}
		}
//		fmt.Println(line)
	}
	input = append(input, strings.Join(trimmedFile, "\n"))
}

func main() {
	var ok bool

	input := strings.FieldsFunc(input[0], whitespaceSplitter)

	if headDefs, ok = rcDefs["head defs"].(gomn.Map); !ok {
		kilOcont("head defs not defined")
	}

	if importsMap, ok = headDefs["imports"].(gomn.Map); !ok {
		kilOcont("imports map not found, could be a non-problem")
	} else if importDefs, ok = importsMap["defs"].(gomn.Map); !ok {
		kilOcont("imports map found, but no definitions found") 
	}

	//parse the header
	var outputHeader []string
	outputHeader = parseHeader(inputHeader, outputHeader)

	//parse the main script
	var outputMain []string
	outputMain = parse(input, outputMain, defsGlob, false)

	//combine them for output
	output := make([]string, len(outputMain)+len(outputHeader))
	copy(output, outputHeader)
	copy(output[len(outputHeader):], outputMain)

	//print it
	//  (for testing, will be changed to write to file)
	for i, chunk := range output {
		if len(splitters)-1 < i {
			splitters = append(splitters, "\n")
		}
		fmt.Print(chunk + splitters[i])
	}//	fmt.Printf("new: %#v\n", output)
}


func importParser(old []string, defs gomn.Map) []string {
	var out []string
	return out
}

func appOut(old []string, cond bool, newVal string, oldVal string) []string {
	if cond {
		return append(old, newVal) 
	} else {
		return append(old, oldVal)
	}
	return []string{}
}

func whitespaceSplitter(r rune) bool {
	if r == '"' {
		isString = !isString
		return false
	} else if isString {
		return false
	} else {
		switch r {
		case '\n', ' ':
			splitters = append(splitters, string(r))
			return true
			break
		case '.':
			splitters = append(splitters, ".")
			return false
		default:
			break
		}
	}
	return false
}

func kilOcont(str string) {
	log.Warn(str)
	if killOnWarn {
		os.Exit(1)
	} else {
		log.Info("continuing anyways")
	}
}
