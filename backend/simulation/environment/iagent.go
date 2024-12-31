package environment

type AgentID string

type Orientation string

const (
	North     Orientation = "N"
	East      Orientation = "E"
	South     Orientation = "S"
	West      Orientation = "W"
	NorthEast Orientation = "NE"
	NorthWest Orientation = "NW"
	SouthEast Orientation = "SE"
	SouthWest Orientation = "SW"
)

type AgentType string

const (
	ExAgent AgentType = "ExAgent"
	Bee     AgentType = "Bee"
	Hornet  AgentType = "Hornet"
)

// IAgent is an interface representing the agent's actions
// limits
type IAgent interface {
	ID() AgentID
	Position() *Position
	Orientation() Orientation
	GetSyncChan() chan bool
	ToJsonObj() interface{}
	Start()
	Percept()
	Deliberate()
	Act()
	Type() AgentType
}
