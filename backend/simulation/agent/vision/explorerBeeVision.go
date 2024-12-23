package vision

import (
	"math"
	"sort"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

// For now we decide that the vision of a bee is an equilateral triangle
// We decide that the returned list is sorted by proximity to the agent
func ExplorerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []SeenElem {
	// Height of the triangle
	distance := utils.GetBeeAgentVisionRange()
	// Side of the triangle - formula for an equilateral triangle
	sideSize := int((2*(distance))/math.Sqrt(3)) + 1
	// Vision triangle coordinates
	topCorner, leftCorner, rightCorner := getTriangleCoordinates(*agt.Position(), int(distance), sideSize, agt.Orientation())

	// To avoid stringy code
	agtX, agtY := agt.Position().X, agt.Position().Y

	// Contains all the elements seen by the agent
	seenElems := make([]SeenElem, 0)

	addElemToList := func(x, y int) {
		if env.IsValidPosition(x, y) {
			if pointIsInTriangle(x, y, topCorner, leftCorner, rightCorner) {
				seenElems = append(seenElems, *NewSeenElem(envpkg.Position{X: x, Y: y}, env.GetAt(x, y)))
			}
		}
	}

	for x := agtX - int(distance/2) + 1; x <= agtX+int(distance/2)+1; x++ {
		for y := agtY - int(distance/2) + 1; y >= agtY+int(distance/2)+1; y++ {
			addElemToList(x, y)
		}
	}

	// Sorting the list by proximity to the agent
	sort.Slice(seenElems, func(i, j int) bool {
		return agt.Position().DistanceFrom(seenElems[i].Pos) < agt.Position().DistanceFrom(seenElems[j].Pos)
	})

	return seenElems
}
