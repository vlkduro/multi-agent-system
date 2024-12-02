package environment

type ObjectID string

type IObject interface {
	ID() ObjectID
	Position() *Position
	Copy() interface{}
	Become(interface{})
	ToJsonObj() interface{}
	Interact()
}
