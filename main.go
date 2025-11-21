package main

import (
	"os"
	"strings"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	//sorry, but I have a lot of
	// vars to keep track of things 
	//  and for configs
	args = os.Args[1:]
	er = log.New(os.Stderr)
	isString bool
	fileExt string
	execute bool
	killOnWarn bool
	writeFile bool
	printOut bool
	debug bool
	input []string
	inputFile string
	outputFile string
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
	readConf()
	checkArgs()

	//i get an odd bug if not defined this way
	var err error

	if inputFile == "" {
		log.Debug("input file not set")
	} else { log.Debug("input file set") }
	
	log.Debug("reading input file")
	inFile, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	} else { log.Debug("success reading input file") }
	
	log.Debug("creating trimming file")
	strFile := string(inFile)
	var trimmedFile []string
	for _, line := range strings.Split(strFile, "\n") {
		log.Debug(line)
		if trimmedLine := strings.TrimSpace(line); trimmedLine != headEnd && trimmedLine != "" {
			log.Debug("not matched as head end")
			log.Debug("adding trimmed line")
			trimmedFile = append(trimmedFile, trimmedLine)
		} else if trimmedLine != "" {
			log.Debug("matched as head end")
			log.Debug("storing input header separately")
			log.Debug("joining header lines") 
			conHeader := strings.Join(trimmedFile, "\n")

			log.Debug("splitting header into chunks")
			splitHeader := strings.FieldsFunc(conHeader, whitespaceSplitter)

			log.Debug("setting header") 
			inputHeader = splitHeader

			log.Debug("reseting trimmed file")
			trimmedFile = []string{""}
		}
//		fmt.Println(line)
	}

	log.Debug("joining trimmed main script into string")
	joinedTrimMain := strings.Join(trimmedFile, "\n")

	log.Debug("appending trimmed main script string to empty input slice")
	input = append(input, joinedTrimMain)
	
	log.Debug("initialized")
}

func main() {
	log.Debug("splitting item one of input slice and over-writing input var")
	input := strings.FieldsFunc(input[0], whitespaceSplitter)

	log.Debug("parsing header")
	//parse the header
	var outputHeader []string
	outputHeader = parseHeader(inputHeader, outputHeader)

	log.Debug("parsing main script")
	//parse the main script
	var outputMain []string
	outputMain = parse(input, outputMain, defsGlob, false)

	log.Debug("combining header and main into output") 
	//combine them for output
	output := make([]string, len(outputMain)+len(outputHeader))
	copy(output, outputHeader)
	copy(output[len(outputHeader):], outputMain)

	log.Debugf("new: %#v\n", output)

	log.Debug("constructing final file")
	var finalOut string
	for i, chunk := range output {
		log.Debugf("line:  %d", i)
		if len(splitters)-1 < i {
			log.Debug("splitter arr ended but chunks have not, adding newline")
			splitters = append(splitters, "\n")
		}
		
		log.Debug("adding chunk and splitter to final output")
		finalOut += chunk + splitters[i]
	}

	if printOut {
		log.Debug("printing output")
		log.Print(finalOut)
	} else { log.Debug("not printing output") }

	if writeFile {
		log.Debug("writing file")

		if outputFile == "" {
			log.Debug("output file not set")

			log.Debug("getting original file extension")
			orExt := filepath.Ext(inputFile)

			log.Debug("getting original filename and removing extension")
			orName := strings.TrimSuffix(filepath.Base(inputFile), orExt)

			log.Debug("building output filename")
			outputFile = orName + "." + fileExt

			log.Warn("no output file, using input file name:  " + outputFile)
		} else { log.Debug("output file is already set") }
		
		if err := os.WriteFile(outputFile, []byte(finalOut), 0644); err != nil {
			log.Fatalf("failed to write to file:  %v", err)
		} else { log.Debug("success writing to file") }

	} else { log.Debug("not writing to file") }
	
	log.Debug("completed") 
}
