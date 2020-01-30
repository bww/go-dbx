package pql

import (
	"unicode"
	"unicode/utf8"
)

const (
	eof     rune = -1
	invalid rune = utf8.RuneError
)

type Scanner struct {
	text   string
	index  int
	width  int
	nchar  rune
	nwidth int
}

func NewScanner(t string) *Scanner {
	return &Scanner{text: t}
}

func (s *Scanner) Text() string {
	return s.text[s.index:]
}

func (s *Scanner) Substr(l, u int) string {
	return s.text[l:u]
}

func (s *Scanner) Peek() rune {
	if s.nwidth > 0 {
		return s.nchar
	}
	if s.index >= len(s.text) {
		return eof
	}
	s.nchar, s.nwidth = utf8.DecodeRuneInString(s.text[s.index:])
	return s.nchar
}

func (s *Scanner) Next() rune {
	var r rune
	if s.nwidth > 0 {
		r, s.nchar = s.nchar, 0
		s.width, s.nwidth = s.nwidth, 0
		s.index += s.width
	} else if s.index >= len(s.text) {
		r = eof
	} else {
		r, s.width = utf8.DecodeRuneInString(s.text[s.index:])
		s.index += s.width
	}
	return r
}

func (s *Scanner) Step() *Scanner {
	s.Next() // advance and ignore
	return s
}

func (s *Scanner) Backup() *Scanner {
	w := s.width
	s.index -= w
	s.width = 0 // can only backup once
	return s
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

func (s *Scanner) SkipWhite() *Scanner {
	s.Skip(unicode.IsSpace)
	return s
}
