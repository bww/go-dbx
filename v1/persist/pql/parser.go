package pql

import (
	"strings"
	"unicode"
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
		if e != nil {
			n = append(n, e)
		}
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
			return parseExprList(s)
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

func parseExprList(s *Scanner) (Node, error) {
	sub := make([]Node, 0)
	for {
		s.SkipWhite()
		c := s.Next()
		if c == eof {
			return nil, newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
		}

		e, err := parseExpr(s)
		if err != nil {
			return nil, err
		}
		sub = append(sub, e)

		s.SkipWhite()
		c = s.Next()
		if c == eof {
			return nil, newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
		} else if c == ',' {
			continue
		} else if c == '}' {
			break
		} else {
			return nil, newErr(ErrUnexpectedToken, NewSpan(s.text, s.index, 1))
		}
	}
	return nil, nil
}

func parseExpr(s *Scanner) (Node, error) {
	s.SkipWhite()
	a := s.index

	pfx, name, err := parseQName(s)
	if err != nil {
		return nil, err
	}

	if pfx != "" && name != "" {
		return exprLiteralNode{node: newNode(s.text, a, s.index-a), prefix: pfx, name: name}, nil
	} else if pfx == "" && name != "" {
		return exprLiteralNode{node: newNode(s.text, a, s.index-a), name: name}, nil
	} else if pfx != "" && name == "" {
		return exprMatchNode{node: newNode(s.text, a, s.index-a), prefix: pfx}, nil
	} else { // prefix and name nil
		return exprMatchNode{node: newNode(s.text, a, s.index-a)}, nil
	}
}

func parseQName(s *Scanner) (string, string, error) {
	pfx, err := parseIdent(s)
	if err != nil {
		return "", "", err
	}

	s.SkipWhite()
	c := s.Next()
	if c == eof {
		return "", "", newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
	} else if c != '.' {
		s.Backup()
		return "", pfx, nil
	}

	c = s.Next()
	if c == eof {
		return "", "", newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
	} else if c == '*' {
		return pfx, "", nil
	} else {
		s.Backup()
	}

	name, err := parseIdent(s)
	if err != nil {
		return "", "", err
	}

	return pfx, name, nil
}

func parseIdent(s *Scanner) (string, error) {
	a := s.index
	c := s.Next()
	if c != '_' && !unicode.IsLetter(c) && !unicode.IsDigit(c) {
		return "", newErr(ErrInvalidIdent, NewSpan(s.text, a, s.index-a))
	}
	for c = s.Next(); c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c); {
		c = s.Next()
	}
	s.Backup()
	return s.text[a:s.index], nil
}
