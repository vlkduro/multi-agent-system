package environment

import (
	agt "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent"
)

type Position struct {
	maxX int
	maxY int
	X    int `json:"x"`
	Y    int `json:"y"`
}

func NewPosition(x int, y int, maxX int, maxY int) *Position {
	return &Position{maxX: maxX, maxY: maxY, X: x, Y: y}
}

func (p *Position) GoUp(ag agt.Agent) *Position {
	p.Y = (p.Y + 1) * ag.Speed
	return p
}

func (p *Position) GoDown(ag agt.Agent) *Position {
	p.Y = (p.Y - 1) * ag.Speed
	return p
}

func (p *Position) GoLeft(ag agt.Agent) *Position {
	p.X = (p.X - 1) * ag.Speed
	return p
}

func (p *Position) GoRight(ag agt.Agent) *Position {
	p.X = (p.X + 1) * ag.Speed
	return p
}

func (p Position) Copy() *Position {
	return &Position{X: p.X, Y: p.Y}
}

func (p Position) Equal(other *Position) bool {
	return p.X == other.X && p.Y == other.Y
}
