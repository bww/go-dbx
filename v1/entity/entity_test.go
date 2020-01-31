package entity

type embedEntity struct {
	B string `db:"y"`
}

type testEntity struct {
	embedEntity
	A string `db:"z,pk"`
	D int    // ignored
}

type multiPKEntity struct {
	embedEntity
	A string `db:"z,pk"`
	C string `db:"x,pk"`
}
