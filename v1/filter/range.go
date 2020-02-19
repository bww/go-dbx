package filter

var ZeroRange = Range{}

type Range struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}
