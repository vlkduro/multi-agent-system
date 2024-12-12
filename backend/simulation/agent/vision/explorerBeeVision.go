package vision

import (
	"math"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// For now we decide that the vision of a bee is an equilateral triangle
func ExplorerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []SeenElem {
	// Height of the triangle
	distance := 6.0
	// Side of the triangle - formula for an equilateral triangle
	sideSize := int((2*(distance))/math.Sqrt(3)) + 1
	// Vision triangle coordinates
	topCorner, leftCorner, rightCorner := getTriangleCoordinates(*agt.Position(), int(distance), sideSize, agt.Orientation())

	// To avoid stringy code
	agtX, agtY := agt.Position().X, agt.Position().Y

	// Contains all the elements seen by the agent
	seenElems := make([]SeenElem, 0)

	for x := agtX - int(distance/2) + 1; x <= agtX+int(distance/2)+1; x++ {
		for y := agtY - int(distance/2) + 1; y <= agtY+int(distance/2)+1; y++ {
			if env.IsValidPosition(x, y) {
				if pointIsInTriangle(x, y, topCorner, leftCorner, rightCorner) {
					seenElems = append(seenElems, *NewSeenElem(envpkg.Position{X: x, Y: y}, env.GetAt(x, y)))
				}
			}
		}
	}
	return seenElems
}
