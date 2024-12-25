package vision

import (
	"fmt"
	"math"
	"sort"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

// For now we decide that the vision of a bee is an equilateral triangle
// We decide that the returned list is sorted by proximity to the agent
func ExplorerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []*SeenElem {
	// Height of the triangle
	distance := utils.GetBeeAgentVisionRange()
	// Side of the triangle - formula for an equilateral triangle
	sideSize := (2*(distance))/math.Sqrt(3) + 1
	// Vision triangle coordinates
	topCornerX, topCornerY, leftCornerX, leftCornerY, rightCornerX, rightCornerY := getTriangleCoordinates(*agt.Position(), distance, sideSize, agt.Orientation())
	// Getting the bounding box of the triangle
	minX := utils.Min(topCornerX, utils.Min(leftCornerX, rightCornerX))
	maxX := utils.Max(topCornerX, utils.Max(leftCornerX, rightCornerX))
	minY := utils.Min(topCornerY, utils.Min(leftCornerY, rightCornerY))
	maxY := utils.Max(topCornerY, utils.Max(leftCornerY, rightCornerY))

	// Contains all the elements seen by the agent
	seenElems := make([]*SeenElem, 0)

	addElemToList := func(x, y int) {
		if env.IsValidPosition(x, y) {
			if pointIsInTriangle(float64(x), float64(y), topCornerX, topCornerY, leftCornerX, leftCornerY, rightCornerX, rightCornerY) {
				seenElems = append(seenElems, NewSeenElem(envpkg.Position{X: x, Y: y}, env.GetAt(x, y)))
				if env.GetAt(x, y) != nil {
					fmt.Printf("[%s] Found something at (%d %d) : %v\n", agt.ID(), x, y, env.GetAt(x, y))
				}
			}
		}
	}

	for x := utils.Round(minX) + 1; x <= utils.Round(maxX)+1; x++ {
		for y := utils.Round(minY); y <= utils.Round(maxY)+1; y++ {
			addElemToList(x, y)
		}
	}

	// Sorting the list by proximity to the agent
	sort.Slice(seenElems, func(i, j int) bool {
		return agt.Position().DistanceFrom(seenElems[i].Pos) < agt.Position().DistanceFrom(seenElems[j].Pos)
	})

	return seenElems
}
