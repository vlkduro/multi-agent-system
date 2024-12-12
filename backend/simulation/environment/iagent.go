package environment

type AgentID string

type Orientation string

const (
	North Orientation = "N"
	East  Orientation = "E"
	South Orientation = "S"
	West  Orientation = "W"
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
}
