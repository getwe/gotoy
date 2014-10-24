package gracetext

import (
	"fmt"
)

type Line struct {
	StrList []string
}

func NewLine(txt string) *Line {
	l := &Line{}
	l.StrList = []string{txt}
	return l
}

func NewLineList(strlist []string) *Line {
	l := &Line{}
	l.StrList = strlist
	return l
}

func (this *Line) Print() {
	for _, txt := range this.StrList {
		fmt.Println(txt)
	}
}

type TextProcess interface {
	Do(*Line) *Line
}
