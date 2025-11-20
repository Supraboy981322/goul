package main

//import "github.com/charmbracelet/log"

var (
	isImp bool
)

func getHeader(in []string) []string {
	var header []string
	for i, val := range in {
		if val == headEnd {
			header = in[:i]
			input = in[i+1:]
		}
	}

	return header
}
	
func parseHeader(in []string, out []string) []string {
	for _, chunk := range in {
		if newChunk, ok := headDefs[chunk].(string); ok {
			out = append(out, newChunk)
		} else if impDefs, ok := importsMap[chunk].(string); ok {
			isImp = !isImp
			out = append(out, impDefs)
		} else if imp, ok := importDefs[chunk].(string); ok && isImp {
			out = append(out, imp)
		} else {
			out = append(out, chunk)
		}
	}

	return out
}
