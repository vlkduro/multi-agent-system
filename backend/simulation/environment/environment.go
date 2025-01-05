package environment

import (
	"fmt"
	"sync"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

var mapDimension int

type Environment struct {
	sync.Mutex
	agts []IAgent
	objs []IObject
	grid [][]interface{}
}

type EnvironmentJson struct {
	MapDimension int             `json:"mapDimension"`
	Grid         [][]interface{} `json:"grid"`
}

func NewEnvironment(agts []IAgent, objs []IObject) *Environment {
	mapDimension = utils.GetMapDimension()
	grid := make([][]interface{}, mapDimension)
	for i := range grid {
		grid[i] = make([]interface{}, mapDimension)
	}
	return &Environment{agts: agts, objs: objs, grid: grid}
}

func (env *Environment) GetAt(x int, y int) interface{} {
	if x < 0 || y < 0 || x >= mapDimension || y >= mapDimension {
		return nil
	}
	elem := env.grid[x][y]
	if elem == nil {
		return nil
	}
	switch obj := elem.(type) {
	case IAgent:
		return obj
	case IObject:
		return obj
	default:
		return obj
	}
}

func (env *Environment) GetMap() [][]interface{} {
	return env.grid
}

func (env *Environment) GetMapDimension() int {
	return mapDimension
}

func (env *Environment) IsValidPosition(x int, y int) bool {
	return x >= 0 && y >= 0 && x < mapDimension && y < mapDimension
}

// Returns true is added on grid, false if not
func (env *Environment) AddAgent(agt IAgent) bool {
	pos := agt.Position()
	env.agts = append(env.agts, agt)
	if pos != nil && env.GetAt(pos.X, pos.Y) == nil {
		env.grid[pos.X][pos.Y] = agt
		return false
	}
	return true
}

func (env *Environment) RemoveAgent(agt IAgent) {
	for i, a := range env.agts {
		if a.ID() == agt.ID() {
			env.grid[agt.Position().X][agt.Position().Y] = nil
			env.agts = append(env.agts[:i], env.agts[i+1:]...)
			break
		}
	}
}

func (env *Environment) RemoveObject(obj IObject) {
	for i, a := range env.objs {
		if a.ID() == obj.ID() {
			env.objs = append(env.objs[:i], env.objs[i+1:]...)
			env.grid[obj.Position().X][obj.Position().Y] = nil
			break
		}
	}
}

func (env *Environment) AddObject(obj IObject) bool {
	pos := obj.Position()
	if env.GetAt(pos.X, pos.Y) != nil {
		return false
	}
	env.objs = append(env.objs, obj)
	env.grid[pos.X][pos.Y] = obj
	return true
}

// https://web.archive.org/web/20171022224528/http://www.policyalmanac.org:80/games/aStarTutorial.htm
func (env *Environment) PathFinding(start *Position, end *Position, numberMoves int) []*Position {
	if start == nil || end == nil {
		return nil
	}
	if start.Equal(end) {
		return make([]*Position, 0)
	}
	type node struct {
		position  *Position
		parent    *node
		cost      float64
		heuristic float64
	}

	if !env.IsValidPosition(start.X, start.Y) || !env.IsValidPosition(end.X, end.Y) {
		return nil
	}

	openList := []*node{}
	closedList := map[string]bool{}
	posToString := func(pos *Position) string {
		return fmt.Sprintf("%d %d", pos.X, pos.Y)
	}

	startNode := &node{position: start.Copy(), cost: 0, heuristic: start.ManhattanDistance(end)}
	openList = append(openList, startNode)

	// We allow pathfinding to last 3 times the number of moves
	cpt := 0
	for len(openList) > 0 && cpt < numberMoves*3 {
		//fmt.Printf("%d ", cpt)
		currentNode := openList[0]
		currentIndex := 0
		for index, node := range openList {
			if node.cost+node.heuristic < currentNode.cost+currentNode.heuristic {
				currentNode = node
				currentIndex = index
			}
		}

		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedList[posToString(currentNode.position)] = true

		if currentNode.position.Equal(end) {
			path := []*Position{}
			for currentNode != nil {
				path = append([]*Position{currentNode.position.Copy()}, path...)
				currentNode = currentNode.parent
			}
			return path
		}

		for _, neighbor := range currentNode.position.GetNeighbours(1) {
			if !env.IsValidPosition(neighbor.X, neighbor.Y) || closedList[posToString(neighbor)] {
				continue
			}

			if _, ok := env.GetAt(neighbor.X, neighbor.Y).(IAgent); ok {
				continue
			}

			cost := currentNode.cost + 1
			heuristic := neighbor.ManhattanDistance(end)
			neighborNode := &node{position: neighbor.Copy(), parent: currentNode, cost: cost, heuristic: heuristic}
			openList = append(openList, neighborNode)
		}
		cpt++
	}
	//fmt.Println("")

	path := []*Position{}
	var currentNode *node = nil
	if len(openList) > 0 {
		currentNode = openList[len(openList)-1]
	}
	for currentNode != nil {
		path = append([]*Position{currentNode.position.Copy()}, path...)
		currentNode = currentNode.parent
	}
	return path
}

func (env *Environment) ToJsonObj() interface{} {
	grid := make([][]interface{}, mapDimension)
	for i := range grid {
		grid[i] = make([]interface{}, mapDimension)
	}

	for x := range env.grid {
		for y := range env.grid[x] {
			gridObj := env.grid[x][y]
			if gridObj != nil {
				switch obj := gridObj.(type) {
				case IAgent:
					grid[x][y] = obj.ToJsonObj()
				case IObject:
					grid[x][y] = obj.ToJsonObj()
				}
			}
		}
	}

	return EnvironmentJson{MapDimension: mapDimension, Grid: grid}
}

func (env *Environment) GetHive() IObject {
	return env.objs[0]
}

func (env *Environment) GetNumberAgent() int {
	return len(env.agts)
}
