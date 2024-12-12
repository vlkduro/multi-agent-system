package vision

import (
	"math"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

type SeenElem struct {
	Pos  envpkg.Position `json:"position"`
	Elem interface{}     `json:"elem"`
}

func NewSeenElem(pos envpkg.Position, elem interface{}) *SeenElem {
	return &SeenElem{Pos: pos, Elem: elem}
}

type VisionFunc func(agt envpkg.IAgent, env *envpkg.Environment) []SeenElem

// To understanc the startPt, leftCorner, rightCorner, and orientation parameters, we need to look at the triangle from the base up, startPt being the top of the triangle
// Source : Alexandre (donc pas sur que ca soit optimal mais en vrai la complexite est O(1) donc ca va)
func getTriangleCoordinates(startPt envpkg.Position, height int, oppositeBaseSize int, orientation envpkg.Orientation) (topCorner, leftCorner, rightCorner envpkg.Position) {
	topCorner = *startPt.Copy()
	switch orientation {
	case envpkg.North:
		rightCorner = envpkg.Position{X: startPt.X - oppositeBaseSize/2, Y: startPt.Y - height}
		leftCorner = envpkg.Position{X: startPt.X + oppositeBaseSize/2, Y: startPt.Y - height}
	case envpkg.South:
		leftCorner = envpkg.Position{X: startPt.X - oppositeBaseSize/2, Y: startPt.Y + height}
		rightCorner = envpkg.Position{X: startPt.X + oppositeBaseSize/2, Y: startPt.Y + height}
	case envpkg.East:
		rightCorner = envpkg.Position{X: startPt.X + height, Y: startPt.Y - oppositeBaseSize/2}
		leftCorner = envpkg.Position{X: startPt.X + height, Y: startPt.Y + oppositeBaseSize/2}
	case envpkg.West:
		leftCorner = envpkg.Position{X: startPt.X - height, Y: startPt.Y - oppositeBaseSize/2}
		rightCorner = envpkg.Position{X: startPt.X - height, Y: startPt.Y + oppositeBaseSize/2}
	}
	return
}

// Source : https://www.geeksforgeeks.org/check-whether-a-given-point-lies-inside-a-triangle-or-not/
func pointIsInTriangle(x, y int, topCorner, leftCorner, rightCorner envpkg.Position) bool {

	calculateArea := func(x1, y1, x2, y2, x3, y3 int) float64 {
		return math.Abs(float64(x1*(y2-y3)+x2*(y3-y1)+x3*(y1-y2)) / 2.0)
	}

	area := calculateArea(topCorner.X, topCorner.Y, leftCorner.X, leftCorner.Y, rightCorner.X, rightCorner.Y)
	area1 := calculateArea(x, y, leftCorner.X, leftCorner.Y, rightCorner.X, rightCorner.Y)
	area2 := calculateArea(topCorner.X, topCorner.Y, x, y, rightCorner.X, rightCorner.Y)
	area3 := calculateArea(topCorner.X, topCorner.Y, leftCorner.X, leftCorner.Y, x, y)

	return area == area1+area2+area3

}
