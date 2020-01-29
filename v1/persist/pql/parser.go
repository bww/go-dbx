package pql

import (
	"strings"
)

func Parse(t string) (*Program, error) {
	s := NewScanner(t)
	n := make([]Node, 0)
	for {
		e, err := parseNode(s)
		if err == EOF {
			break
		} else if err != nil {
			return nil, err
		}
		n = append(n, e)
	}
	return &Program{sub: n}, nil
}

func parseNode(s *Scanner) (Node, error) {
	for {
		c := s.Next()
		if c == eof {
			break
		}
		switch c {
		case '{':
			return parseMeta(s)
		default:
			return parseLiteral(s)
		}
	}
	return nil, EOF
}

func parseLiteral(s *Scanner) (Node, error) {
	var err error
	b := &strings.Builder{}
	a := s.index
	for {
		c := s.Next()
		if c == eof {
			break
		} else if c == '{' {
			s.Backup()
			break // found meta
		} else if c == '\\' {
			c, err = parseEscape(s)
			if err != nil {
				return nil, err
			}
		}
		b.WriteRune(c)
	}
	return literalNode{
		node: newNode(s.text, a, s.index-a),
		text: b.String(),
	}, nil
}

func parseEscape(s *Scanner) (rune, error) {
	a := s.index
	switch c := s.Next(); c {
	case '\\':
		return '\\', nil
	case '{':
		return '{', nil
	default:
		return eof, newErr(ErrInvalidEscape, NewSpan(s.text, a, s.index-a))
	}
}

func parseMeta(s *Scanner) (Node, error) {
	return nil, nil
}
