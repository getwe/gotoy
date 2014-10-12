package main

import (
	"bufio"
	. "github.com/getwe/gotoy/gracestdout/gracetext"
	"os"
)

func main() {

	handler := make([]TextProcess, 0)

	handler = append(handler, new(JsonTextProcess))
	handler = append(handler, new(FmtTextProcess))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		for _, h := range handler {
			if h.Do(text) {
				break
			}
		}
	}
}
