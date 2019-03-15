package prompt

type Coord struct {
	X Column
	Y Row
}

func (c Coord) Add(b Coord) Coord {
	return Coord{c.X + b.X, c.Y + b.Y}
}

func (c Coord) Diff(b Coord) Coord {
	return Coord{c.X - b.X, c.Y - b.Y}
}
