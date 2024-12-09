package environment

type AgentID string

// IAgent is an interface representing the agent's actions
// limits
type IAgent interface {
	ID() AgentID
	Position() *Position
	GetSyncChan() chan bool
	ToJsonObj() interface{}
	Start()
	Percept()
	Deliberate()
	Act()
}
