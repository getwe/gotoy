package main

import (
	"bufio"
	. "github.com/getwe/gotoy/gracestdout/gracetext"
	"os"
)

func main() {

	handler := make([]TextProcess, 0)

	handler = append(handler, new(JsonTextProcess))
	handler = append(handler, new(JsonColorProcess))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()

		line := NewLine(text)

		for _, h := range handler {
			line = h.Do(line)
		}

		line.Print()
	}
}
