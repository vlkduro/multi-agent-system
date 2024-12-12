package environment

// Coordinate system based on the top left corner
type Position struct {
	maxX int
	maxY int
	X    int `json:"x"`
	Y    int `json:"y"`
}

func NewPosition(x int, y int, maxX int, maxY int) *Position {
	return &Position{maxX: maxX, maxY: maxY, X: x, Y: y}
}

func (p *Position) move(grid [][]interface{}, newX int, newY int) *Position {
	if grid[newX][newY] != nil {
		return p
	}
	grid[p.X][p.Y], grid[newX][newY] = nil, grid[p.X][p.Y]
	p.X = newX
	p.Y = newY
	return p
}

func (p *Position) GoUp(grid [][]interface{}) *Position {
	if p.Y == 0 {
		return p
	}
	return p.move(grid, p.X, p.Y-1)
}

func (p *Position) GoDown(grid [][]interface{}) *Position {
	if p.Y == p.maxY-1 {
		return p
	}
	return p.move(grid, p.X, p.Y+1)
}

func (p *Position) GoLeft(grid [][]interface{}) *Position {
	if p.X == 0 {
		return p
	}
	return p.move(grid, p.X-1, p.Y)
}

func (p *Position) GoRight(grid [][]interface{}) *Position {
	if p.X == p.maxX-1 {
		return p
	}
	return p.move(grid, p.X+1, p.Y)
}

func (p Position) Copy() *Position {
	return &Position{X: p.X, Y: p.Y, maxX: p.maxX, maxY: p.maxY}
}

func (p Position) Equal(other *Position) bool {
	return p.X == other.X && p.Y == other.Y
}
