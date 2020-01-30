package pql

import (
	"fmt"
	"strconv"
)

type Span struct {
	text   string
	offset int
	length int
}

func NewSpan(t string, o, l int) Span {
	return Span{t, o, l}
}

func (s Span) Excerpt() string {
	max := len(s.text)
	return s.text[imax(0, imin(max, s.offset)):imin(max, s.offset+s.length)]
}

func (s Span) Describe() string {
	return fmt.Sprintf("[%d+%d] %s", s.offset, s.length, strconv.Quote(s.Excerpt()))
}

func (s Span) String() string {
	return strconv.Quote(s.Excerpt())
}

func Encompass(a ...Span) Span {
	var t string
	min, max := 0, 0
	for i, e := range a {
		if i == 0 {
			min, max = e.offset, e.offset+e.length
			t = e.text
		} else {
			if e.offset < min {
				min = e.offset
			}
			if e.offset+e.length > max {
				max = e.offset + e.length
			}
		}
	}
	return Span{t, min, max - min}
}

func imin(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func imax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
