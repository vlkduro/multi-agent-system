package environment

type AgentID string

// IAgent is an interface representing the agent's actions
// limits
type IAgent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	GetSyncChan() chan bool
	ToJsonObj() interface{}
}
