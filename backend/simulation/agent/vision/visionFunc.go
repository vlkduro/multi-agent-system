package vision

import (
	"math"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

type SeenElem struct {
	Pos  envpkg.Position `json:"position"`
	Elem interface{}     `json:"elem"`
}

func NewSeenElem(pos envpkg.Position, elem interface{}) *SeenElem {
	return &SeenElem{Pos: pos, Elem: elem}
}

type VisionFunc func(agt envpkg.IAgent, env *envpkg.Environment) []*SeenElem

// To understand the startPt, leftCorner, rightCorner, and orientation parameters, we need to look at the triangle from the base up, startPt being the top of the triangle
// Source : Alexandre (donc pas sur que ca soit optimal mais en vrai la complexite est O(1) donc ca va)
// https://www.youtube.com/watch?v=PSlWb90JJx4 - Je suis désolé mais j'ai pas trouvé beaucoup plus simple
func getTriangleCoordinates(startPt envpkg.Position, height float64, oppositeBaseSize float64, orientation envpkg.Orientation) (topCornerX, topCornerY, leftCornerX, leftCornerY, rightCornerX, rightCornerY float64) {
	topCorner := *startPt.Copy()
	topCornerX = float64(topCorner.X)
	topCornerY = float64(topCorner.Y)
	// We imagine two circles, with the following radiuses :
	// - C1 : The first circle has a radius of height and is centered on the topCorner
	// - C2 : The second circle has a diameter(i.e radius) of oppositeBaseSize(i.e oppositeBaseSize/2) and is centered on the projection
	//        of topCorner on C1
	// Combining the two circles gives us the leftCorner and rightCorner points
	// First circle : C1 = (x - a)² + (y - b)² = height²
	// Second circle : C2 = (x - p)² + (y - q)² = oppositeBaseSize²
	// We get the tangent of the center of C2 (point on C1)
	// The points of intersection between C2 and the tangent of C1 are the leftCorner and rightCorner
	// This allows for maintining precision over the height of the triangle and the opposite base size
	c2radius := oppositeBaseSize / 2
	c2x := topCornerX
	c2y := topCornerY
	var rcxCoef, rcyCoef, lcxCoef, lcyCoef float64
	switch orientation {
	case envpkg.North:
		c2x += height * math.Cos(math.Pi/2)
		c2y += height * math.Sin(math.Pi/2)
		rcxCoef = math.Cos(math.Pi)
		rcyCoef = math.Sin(0)
		lcxCoef = math.Cos(0)
		lcyCoef = math.Sin(0)
	case envpkg.South:
		c2x += height * math.Cos(math.Pi/2)
		c2y += height * math.Sin((3*math.Pi)/2)
		rcxCoef = math.Cos(0)
		rcyCoef = math.Sin(0)
		lcxCoef = math.Cos(math.Pi)
		lcyCoef = math.Sin(0)
	case envpkg.East:
		c2x += height * math.Cos(0)
		c2y += height * math.Sin(0)
		rcxCoef = math.Cos(math.Pi / 2)
		rcyCoef = math.Sin(math.Pi / 2)
		lcxCoef = math.Cos(math.Pi / 2)
		lcyCoef = math.Sin((3 * math.Pi) / 2)
	case envpkg.West:
		c2x += height * math.Cos(math.Pi)
		c2y += height * math.Sin(0)
		rcxCoef = math.Cos(math.Pi / 2)
		rcyCoef = math.Sin((3 * math.Pi) / 2)
		lcxCoef = math.Cos(math.Pi / 2)
		lcyCoef = math.Sin(math.Pi / 2)
	case envpkg.NorthEast:
		c2x += height * math.Cos(math.Pi/4)
		c2y += height * math.Sin(math.Pi/4)
		rcxCoef = math.Cos((5 * math.Pi) / 4)
		rcyCoef = math.Sin(math.Pi / 4)
		lcxCoef = math.Cos(math.Pi / 4)
		lcyCoef = math.Sin((5 * math.Pi) / 4)
	case envpkg.NorthWest:
		c2x += height * math.Cos((5*math.Pi)/4)
		c2y += height * math.Sin(math.Pi/4)
		rcxCoef = math.Cos((5 * math.Pi) / 4)
		rcyCoef = math.Sin((5 * math.Pi) / 4)
		lcxCoef = math.Cos(math.Pi / 4)
		lcyCoef = math.Sin(math.Pi / 4)
	case envpkg.SouthEast:
		c2x += height * math.Cos((5*math.Pi)/4)
		c2y += height * math.Sin(math.Pi/4)
		rcxCoef = math.Cos(math.Pi / 4)
		rcyCoef = math.Sin(math.Pi / 4)
		lcxCoef = math.Cos((5 * math.Pi) / 4)
		lcyCoef = math.Sin((5 * math.Pi) / 4)
	case envpkg.SouthWest:
		c2x += height * math.Cos((5*math.Pi)/4)
		c2y += height * math.Sin((5*math.Pi)/4)
		rcxCoef = math.Cos(math.Pi / 4)
		rcyCoef = math.Sin((5 * math.Pi) / 4)
		lcxCoef = math.Cos((5 * math.Pi) / 4)
		lcyCoef = math.Sin(math.Pi / 4)
	}

	rightCornerX = rcxCoef*c2radius + c2x
	rightCornerY = rcyCoef*c2radius + c2y
	leftCornerX = lcxCoef*c2radius + c2x
	leftCornerY = lcyCoef*c2radius + c2y
	return
}

// Source : https://www.geeksforgeeks.org/check-whether-a-given-point-lies-inside-a-triangle-or-not/
func pointIsInTriangle(x, y, topCornerX, topCornerY, leftCornerX, leftCornerY, rightCornerX, rightCornerY float64) bool {

	calculateArea := func(x1, y1, x2, y2, x3, y3 float64) float64 {
		return math.Abs(float64(x1*(y2-y3)+x2*(y3-y1)+x3*(y1-y2)) / 2.0)
	}

	area := calculateArea(topCornerX, topCornerY, leftCornerX, leftCornerY, rightCornerX, rightCornerY)
	area1 := calculateArea(x, y, leftCornerX, leftCornerY, rightCornerX, rightCornerY)
	area2 := calculateArea(topCornerX, topCornerY, x, y, rightCornerX, rightCornerY)
	area3 := calculateArea(topCornerX, topCornerY, leftCornerX, leftCornerY, x, y)

	return utils.Round(area) == utils.Round(area1+area2+area3)

}
