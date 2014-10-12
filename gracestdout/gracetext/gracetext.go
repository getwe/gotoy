package gracetext

import (
	"bytes"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
)

type TextProcess interface {
	Do(string) bool
}

type JsonTextProcess struct {
}

func (this JsonTextProcess) Do(txt string) bool {
	jsonObj, err := simplejson.NewJson([]byte(txt))
	if err != nil {
		return false
	}
	buf, err := jsonObj.EncodePretty()
	if err != nil {
		return false
	}

	fmt.Println(bytes.NewBuffer(buf).String())

	return true
}

type FmtTextProcess struct {
}

func (this FmtTextProcess) Do(txt string) bool {
	fmt.Println(txt)
	return true
}
