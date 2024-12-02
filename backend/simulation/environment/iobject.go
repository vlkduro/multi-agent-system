package environment

type Object interface {
	Interact()
	Copy() Object
	Become(Object)
	ToJsonObj() interface{}
}
