package pql

import (
	"strings"
	"unicode"
)

const wildcard = "*"

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
	return &Program{
		node: newNode(t, 0, len(t)),
		sub:  n,
	}, nil
}

func parseNode(s *Scanner) (Node, error) {
outer:
	for {
		switch c := s.Peek(); c {
		case eof:
			break outer
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
outer:
	for {
		switch c := s.Next(); c {
		case eof:
			break outer
		case '{':
			s.Backup()
			break outer
		case '\\':
			c, err = parseEscape(s)
			if err != nil {
				return nil, err
			}
			fallthrough
		default:
			b.WriteRune(c)
		}
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
	switch c := s.Next(); c {
	case eof:
		return nil, newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
	case '{':
		return parseExprList(s)
	default:
		return nil, newErr(ErrUnexpectedToken, NewSpan(s.text, s.index, 1))
	}
}

func parseExprList(s *Scanner) (Node, error) {
	a, z := s.index, s.index
	sub := make([]Node, 0)
outer:
	for {
		e, err := parseExpr(s.SkipWhite())
		if err != nil {
			return nil, err
		}
		sub = append(sub, e)
		z = s.index
		switch c := s.SkipWhite().Next(); c {
		case eof:
			return nil, newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
		case ',':
			continue outer
		case '}':
			break outer
		default:
			return nil, newErr(ErrUnexpectedToken, NewSpan(s.text, s.index, 1))
		}
	}
	return exprListNode{
		node: newNode(s.text, a, z-a),
		sub:  sub,
	}, nil
}

func parseExpr(s *Scanner) (Node, error) {
	a := s.index

	names, err := parseQName(s)
	if err != nil {
		return nil, err
	}

	if l := len(names); l == 1 {
		if names[0] == wildcard {
			return exprMatchNode{node: newNode(s.text, a, s.index-a)}, nil
		} else {
			return exprLiteralNode{node: newNode(s.text, a, s.index-a), name: names[0]}, nil
		}
	} else if l == 2 {
		if names[0] == wildcard {
			return nil, newErr(ErrInvalidQName, NewSpan(s.text, a, s.index-a))
		} else if names[1] == wildcard {
			return exprMatchNode{node: newNode(s.text, a, s.index-a), prefix: names[0]}, nil
		} else {
			return exprLiteralNode{node: newNode(s.text, a, s.index-a), prefix: names[0], name: names[1]}, nil
		}
	}

	return nil, newErr(ErrInvalidQName, NewSpan(s.text, a, s.index-a))
}

func parseQName(s *Scanner) ([]string, error) {
	var names []string
outer:
	for {
		n, err := parseWildcardIdent(s)
		if err != nil {
			return nil, err
		}

		names = append(names, n)

		switch c := s.SkipWhite().Peek(); c {
		case eof:
			return nil, newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
		case '.':
			s.Next()
		default:
			break outer
		}
	}
	return names, nil
}

func parseWildcardIdent(s *Scanner) (string, error) {
	switch c := s.SkipWhite().Next(); c {
	case eof:
		return "", newErr(ErrUnexpectedEOF, NewSpan(s.text, s.index, 0))
	case '*':
		return wildcard, nil
	}
	name, err := parseIdent(s.Backup())
	if err != nil {
		return "", err
	}
	return name, nil
}

func parseIdent(s *Scanner) (string, error) {
	a := s.index
	for {
		c := s.Peek()
		if c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c) {
			s.Next()
		} else {
			break
		}
	}
	if s.index-a < 1 {
		return "", newErr(ErrInvalidIdent, NewSpan(s.text, a, s.index-a))
	}
	return s.text[a:s.index], nil
}
