package environment

type AgentID string

type Agent interface {
	Start()
	Percept()
	Deliberate()
	Act()
	ID() AgentID
	GetSyncChan() chan bool
}
