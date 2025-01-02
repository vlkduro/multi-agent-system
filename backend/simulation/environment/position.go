package environment

import (
	"math"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

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

// Movement can be simulated with by passing a nil grid
func (p *Position) move(grid [][]interface{}, newX int, newY int) bool {
	// Out of bounds
	if newX < 0 || newY < 0 || newX >= p.maxX || newY >= p.maxY {
		return false
	}
	// Collision
	if grid != nil && (grid[newX][newY] != nil) {
		return false
	}
	// Not a simulation
	if grid != nil {
		grid[p.X][p.Y], grid[newX][newY] = nil, grid[p.X][p.Y]
	}
	p.X = newX
	p.Y = newY
	return true
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoNorth(grid [][]interface{}) bool {
	return p.move(grid, p.X, p.Y-1)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoSouth(grid [][]interface{}) bool {
	return p.move(grid, p.X, p.Y+1)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoWest(grid [][]interface{}) bool {
	return p.move(grid, p.X-1, p.Y)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoEast(grid [][]interface{}) bool {
	return p.move(grid, p.X+1, p.Y)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoNorthEast(grid [][]interface{}) bool {
	return p.move(grid, p.X+1, p.Y-1)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoNorthWest(grid [][]interface{}) bool {
	return p.move(grid, p.X-1, p.Y-1)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoSouthEast(grid [][]interface{}) bool {
	return p.move(grid, p.X+1, p.Y+1)
}

// Movement can be simulated with by passing a nil grid
func (p *Position) GoSouthWest(grid [][]interface{}) bool {
	return p.move(grid, p.X-1, p.Y+1)
}

func (p Position) DistanceFrom(p2 *Position) float64 {
	return math.Sqrt(float64((p.X-p2.X)*(p.X-p2.X) + (p.Y-p2.Y)*(p.Y-p2.Y)))
}

func (p Position) Near(p2 *Position, distance int) bool {
	return utils.Round(p.DistanceFrom(p2)) <= distance
}

// Returns the mirrored position of the point according to this point
func (p Position) GetSymmetricOfPoint(origin Position) *Position {
	symX := 2*origin.X - p.X
	symY := 2*origin.Y - p.Y
	return &Position{X: symX, Y: symY, maxX: p.maxX, maxY: p.maxY}
}

func (p Position) Copy() *Position {
	return &Position{X: p.X, Y: p.Y, maxX: p.maxX, maxY: p.maxY}
}

func (p Position) Equal(other *Position) bool {
	if other == nil {
		return false
	}
	return p.X == other.X && p.Y == other.Y
}

func (p Position) ManhattanDistance(p2 *Position) float64 {
	return math.Abs(float64(p.X-p2.X)) + math.Abs(float64(p.Y-p2.Y))
}

func (p Position) GetNeighbours(distance int) []*Position {
	neighbours := make([]*Position, 0)
	for x := p.X - distance; x <= p.X+distance; x++ {
		for y := p.Y - distance; y <= p.Y+distance; y++ {
			if x == p.X && y == p.Y {
				continue
			}
			if x < 0 || y < 0 || x >= p.maxX || y >= p.maxY {
				continue
			}
			neighbours = append(neighbours, NewPosition(x, y, p.maxX, p.maxY))
		}
	}
	return neighbours
	// return []*Position{
	// 	// NW N NE
	// 	NewPosition(p.X-1, p.Y-1, p.maxX, p.maxY),
	// 	NewPosition(p.X, p.Y-1, p.maxX, p.maxY),
	// 	NewPosition(p.X+1, p.Y-1, p.maxX, p.maxY),
	// 	// W E
	// 	NewPosition(p.X-1, p.Y, p.maxX, p.maxY),
	// 	NewPosition(p.X+1, p.Y, p.maxX, p.maxY),
	// 	// SW S SE
	// 	NewPosition(p.X-1, p.Y+1, p.maxX, p.maxY),
	// 	NewPosition(p.X, p.Y+1, p.maxX, p.maxY),
	// 	NewPosition(p.X+1, p.Y+1, p.maxX, p.maxY),
	// }
}
