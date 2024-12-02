package environment

type Position struct {
	maxX int
	maxY int
	X    int `json:"x"`
	Y    int `json:"y"`
}

func NewPosition(x int, y int, maxX int, maxY int) *Position {
	return &Position{maxX: maxX, maxY: maxY, X: x, Y: y}
}

func (p *Position) GoUp() *Position {
	p.Y++
	return p
}

func (p *Position) GoDown() *Position {
	p.Y--
	return p
}

func (p *Position) GoLeft() *Position {
	p.X--
	return p
}

func (p *Position) GoRight() *Position {
	p.X++
	return p
}

func (p Position) Copy() *Position {
	return &Position{X: p.X, Y: p.Y}
}

func (p Position) Equal(other *Position) bool {
	return p.X == other.X && p.Y == other.Y
}
