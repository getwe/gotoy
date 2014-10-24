package gracetext

import (
	"fmt"
	"github.com/fatih/color"
	"regexp"
)

type JsonColorProcess struct {
}

func (this JsonColorProcess) Do(line *Line) *Line {
	result := make([]string, 0)
	for _, txt := range line.StrList {
		result = append(result, this.strColor(txt))
	}

	return NewLineList(result)
}

func (this JsonColorProcess) strColor(txt string) string {

	keyPattern, _ := regexp.Compile(`(^\s*")(.*?)(":.*$)`)

	if ok := keyPattern.MatchString(txt); ok {
		keyColorFunc := color.New(color.FgYellow).SprintFunc()
		m := keyPattern.FindStringSubmatch(txt)
		return fmt.Sprintf("%s%s%s", m[1], keyColorFunc(m[2]), m[3])
	}

	return txt
}
