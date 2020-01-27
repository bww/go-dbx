package entity

import (
	"reflect"
)

type Values struct {
	Cols []string
	Vals []reflect.Value
}

type Columns struct {
	Cols []string
	Vals []interface{}
}

func (s *Columns) Len() int {
	return len(s.Cols)
}

func (s *Columns) Swap(i, j int) {
	s.Cols[i], s.Cols[j] = s.Cols[j], s.Cols[i]
	s.Vals[i], s.Vals[j] = s.Vals[j], s.Vals[i]
}

func (s *Columns) Less(i, j int) bool {
	return s.Cols[i] < s.Cols[j]
}
