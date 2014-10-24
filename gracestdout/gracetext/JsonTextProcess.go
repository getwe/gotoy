package gracetext

import (
	"bytes"
	simplejson "github.com/bitly/go-simplejson"
	"strings"
)

type JsonTextProcess struct {
}

func (this JsonTextProcess) Do(line *Line) *Line {
	result := make([]string, 0)
	for _, txt := range line.StrList {
		result = append(result, this.str2json(txt)...)
	}

	return NewLineList(result)
}

func (this JsonTextProcess) str2json(txt string) []string {
	jsonObj, err := simplejson.NewJson([]byte(txt))
	if err != nil {
		return []string{txt}
	}
	buf, err := jsonObj.EncodePretty()
	if err != nil {
		return []string{txt}
	}

	multiLine := bytes.NewBuffer(buf).String()
	return strings.Split(multiLine, "\n")
}
