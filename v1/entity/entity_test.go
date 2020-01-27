package entity

type embedEntity struct {
	B string `db:"y"`
}

type testEntity struct {
	embedEntity
	A string `db:"z"`
}
