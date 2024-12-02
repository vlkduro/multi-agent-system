package environment

type IObject interface {
	Interact()
	Copy() IObject
	ChangeTo(IObject)
	ToJsonObj() interface{}
}
