package pql

import (
	"unicode"
	"unicode/utf8"
)

const (
	eof rune = -1
)

type Scanner struct {
	text  string
	index int
	width int // current rune width
}

func NewScanner(t string) *Scanner {
	return &Scanner{text: t}
}

func (s *Scanner) Text() string {
	return s.text[s.index:]
}

func (s *Scanner) Next() rune {
	if s.index >= len(s.text) {
		s.width = 0
		return eof
	} else {
		r, w := utf8.DecodeRuneInString(s.text[s.index:])
		s.index += w
		s.width = w
		return r
	}
}

func (s *Scanner) Backup() int {
	w := s.width
	s.index -= w
	s.width = 0 // can only backup once
	return w
}

func (s *Scanner) Skip(f func(rune) bool) int {
	var i int
	for {
		c := s.Next()
		if c == eof {
			break
		} else if !f(c) {
			s.Backup()
			break
		}
		i++
	}
	return i
}

func (s *Scanner) SkipWhite() int {
	return s.Skip(unicode.IsSpace)
}
