package environment

type Object interface {
	Interact()
	Copy() Object
	ChangeTo(Object)
}
