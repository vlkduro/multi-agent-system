package vision

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// For now we decide that the vision of a bee is an equilateral triangle
// We decide that the returned list is sorted by proximity to the agent
func WorkerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []*SeenElem {
	agt_position := *agt.Position()

	seenElems := make([]*SeenElem, 0)

	addElemToList := func(x, y int) {
		if env.IsValidPosition(x, y) {
			elem := env.GetAt(x, y)
			seenElems = append(seenElems, NewSeenElem(&envpkg.Position{X: x, Y: y}, elem))
		}
	}

	addElemToList(agt_position.X, agt_position.Y)

	return seenElems
}
