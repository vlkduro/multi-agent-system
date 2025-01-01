package vision

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// For now we decide that the vision of a bee is an equilateral triangle
// We decide that the returned list is sorted by proximity to the agent
func WorkerBeeVision(agt envpkg.IAgent, env *envpkg.Environment) []*SeenElem {
	// Contains all the elements seen by the agent
	seenElems := make([]*SeenElem, 0)

	addElemToList := func(x, y int) {
		if env.IsValidPosition(x, y) {
			seenElems = append(seenElems, NewSeenElem(&envpkg.Position{X: x, Y: y}, env.GetAt(x, y)))
		}
	}

	/*fmt.Printf("%v\n", seenElems)
	fmt.Print(env.IsValidPosition(agt.Position().X, agt.Position().Y))
	fmt.Print(env.GetAt(agt.Position().X, agt.Position().Y))

	new_elem := env.GetAt(agt.Position().X, agt.Position().Y)

	seenElems = append(seenElems, new_elem)*/

	addElemToList(agt.Position().X, agt.Position().Y)

	// Sorting the list by proximity to the agent
	/*sort.Slice(seenElems, func(i, j int) bool {
		return agt.Position().DistanceFrom(seenElems[i].Pos) < agt.Position().DistanceFrom(seenElems[j].Pos)
	})*/

	return seenElems
}
