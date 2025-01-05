package vision

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

// For now we decide that the vision of a bee is an equilateral triangle
// We decide that the returned list is sorted by proximity to the agent
func ExplorerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []*SeenElem {
	distance := utils.GetBeeAgentVisionRange()
	return EquilateralTriangleVision(agt, env, distance)
}
