package entity

import (
	"sort"
)

type columnsAndValues struct {
	cols []string
	vals []interface{}
}

func (s *columnsAndValues) Len() int {
	return len(s.cols)
}

func (s *columnsAndValues) Swap(i, j int) {
	s.cols[i], s.cols[j] = s.cols[j], s.cols[i]
	s.vals[i], s.vals[j] = s.vals[j], s.vals[i]
}

func (s *columnsAndValues) Less(i, j int) bool {
	return s.cols[i] < s.cols[j]
}

func sortColumnsAndValues(c []string, v []interface{}) ([]string, []interface{}) {
	cv := &columnsAndValues{c, v}
	sort.Sort(cv)
	return cv.cols, cv.vals
}
