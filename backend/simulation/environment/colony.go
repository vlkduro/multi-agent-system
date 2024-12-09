package environment

type Colony struct {
	bees []IAgent
	// poi = Point of Interest
	chanPoiPostions chan []Position
	// true if is in danger
	inDanger bool
}

func NewColony() *Colony {
	c := &Colony{}
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
