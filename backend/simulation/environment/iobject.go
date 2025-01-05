package environment

type ObjectID string

type ObjectType string

const (
	Flower ObjectType = "Flower"
	Hive   ObjectType = "Hive"
)

type IObject interface {
	ID() ObjectID
	Position() *Position
	Copy() interface{}
	Become(interface{})
	Update()
	ToJsonObj() interface{}
	TypeObject() ObjectType
}
