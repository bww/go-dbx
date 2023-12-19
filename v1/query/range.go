package query

var ZeroRange = Range{}

type Range struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

func (r Range) WithOffset(x int) Range {
	return Range{
		Offset: x,
		Length: r.Length,
	}
}

func (r Range) WithLength(x int) Range {
	return Range{
		Offset: r.Offset,
		Length: x,
	}
}

type Order int

const (
	Ascending Order = iota
	Descending
)
