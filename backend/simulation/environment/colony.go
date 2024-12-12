package environment

// Garder la structure ?
type Colony struct {
	bees []IAgent
	// true if is in danger
	inDanger bool
	// hive
	hive IObject
	// poi = Point of Interest
	chanPoiPostions chan []Position
}

func NewColony(hive IObject) *Colony {
	c := &Colony{hive: hive}
	c.bees = make([]IAgent, 0)
	c.chanPoiPostions = make(chan []Position)
	c.inDanger = false
	return c
}

func (c *Colony) AddBee(bee IAgent) {
	c.bees = append(c.bees, bee)
}

func (c *Colony) GetNumberBees() int {
	return len(c.bees)
}

func (c *Colony) GetHive() IObject {
	return c.hive.Copy().(IObject)
}
