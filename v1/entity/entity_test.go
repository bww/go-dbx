package entity

type embedEntity struct {
	B string `db:"y"`
}

type testEntity struct {
	embedEntity
	A string `db:"z,pk"`
	D int    // ignored
	E int    `db:"e,omitempty"`
}

type multiPKEntity struct {
	embedEntity
	A string `db:"z,pk"`
	C string `db:"x,pk"`
}

type syntheticEntity struct {
	A string `db:"a,pk"`
}

func (s syntheticEntity) AdditionalColumns() *Columns {
	return &Columns{
		Cols: []string{"syn_1"},
		Vals: []interface{}{123},
	}
}
