package environment

type AgentID string

type IAgent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	GetSyncChan() chan bool
	ToJsonObj() interface{}
}
