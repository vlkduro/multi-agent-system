package environment

type AgentID string

type IAgent interface {
	Start()
	Percept(*Environment)
	Deliberate()
	Act(*Environment)
	ID() AgentID
}
