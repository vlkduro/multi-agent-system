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
	return env.grid[x][y]
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

func (env *Environment) AddAgent(agt IAgent) bool {
	pos := agt.Position()
	if env.GetAt(pos.X, pos.Y) == nil {
		//return false
		env.grid[pos.X][pos.Y] = agt
	}
	env.agts = append(env.agts, agt)
	return true
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
func (env *Environment) PathFinding(start *Position, end *Position) []*Position {
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

	// We allow pathfinding to last 50 iterations
	cpt := 50
	for len(openList) > 0 && cpt > 0 {
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

		for _, neighbor := range currentNode.position.GetNeighbours() {
			if !env.IsValidPosition(neighbor.X, neighbor.Y) || closedList[posToString(neighbor)] {
				continue
			}

			if env.GetAt(neighbor.X, neighbor.Y) != nil {
				continue
			}

			cost := currentNode.cost + 1
			heuristic := neighbor.ManhattanDistance(end)
			neighborNode := &node{position: neighbor.Copy(), parent: currentNode, cost: cost, heuristic: heuristic}
			openList = append(openList, neighborNode)
		}
		cpt--
	}

	path := []*Position{}
	currentNode := openList[len(openList)-1]
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
