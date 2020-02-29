package filter

var ZeroRange = Range{}

type Range struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type Order int

const (
	Ascending Order = iota
	Descending
)
