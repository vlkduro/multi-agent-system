package object

import (
	"math/rand/v2"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

// Enum FlowerGender
type FlowerGender string

const (
	Male   FlowerGender = "male"
	Female FlowerGender = "female"
)

type Flower struct {
	envpkg.IObject
	id         envpkg.ObjectID
	pos        *envpkg.Position
	gender     FlowerGender
	pollinated bool
	occupied   bool
	nectar     int // in mg
	maxNectar  int // in mg
}

type FlowerJson struct {
	ID         string          `json:"id"`
	Position   envpkg.Position `json:"position"`
	Gender     FlowerGender    `json:"gender"`
	Pollinated bool            `json:"pollinated"`
	Occupied   bool            `json:"occupied"`
	Nectar     int             `json:"nectar"`
}

func NewFlower(id string, pos *envpkg.Position) *Flower {
	var gender FlowerGender
	if rand.IntN(2) == 0 {
		gender = Male
	} else {
		gender = Female
	}
	maxNectar := utils.GetMaxNectarHeld()
	return &Flower{
		id:         envpkg.ObjectID(id),
		pos:        pos.Copy(),
		gender:     gender,
		pollinated: false,
		occupied:   false,
		nectar:     maxNectar,
		maxNectar:  maxNectar,
	}
}

func (f Flower) ID() envpkg.ObjectID {
	return f.id
}

func (f Flower) Position() *envpkg.Position {
	return f.pos.Copy()
}

func (f Flower) Copy() interface{} {
	return &Flower{
		id:         f.id,
		pos:        f.pos.Copy(),
		gender:     f.gender,
		pollinated: f.pollinated,
		occupied:   f.occupied,
		nectar:     f.nectar,
	}
}

func (f *Flower) Become(f2 interface{}) {
	if f2 == nil {
		return
	}

	f2Flower, ok := f2.(*Flower)
	if !ok {
		return
	}

	f.id = f2Flower.id
	f.pos = f2Flower.pos.Copy()
	f.gender = f2Flower.gender
	f.pollinated = f2Flower.pollinated
	f.occupied = f2Flower.occupied
	f.nectar = f2Flower.nectar
}

func (f *Flower) Update() {
	addedNectar := utils.GetProducedNectarPerTurn()
	f.nectar = (addedNectar + f.nectar)
	if f.nectar > f.maxNectar {
		f.nectar = f.maxNectar
	}
}

func (f Flower) ToJsonObj() interface{} {
	return FlowerJson{
		ID:         string(f.id),
		Position:   *f.pos.Copy(),
		Gender:     f.gender,
		Pollinated: f.pollinated,
		Occupied:   f.occupied,
		Nectar:     f.nectar,
	}
}

func (f *Flower) Occupy() {
	f.occupied = true
}

func (f *Flower) Pollinate() {
	f.pollinated = true
}

func (f *Flower) RetreiveNectar(nectar int) int {
	if f.nectar < nectar {
		nectar = f.nectar
	}
	f.nectar -= nectar
	return nectar
}

func (f *Flower) GetNectar() int {
	return f.nectar
}

func (f Flower) ObjectType() envpkg.ObjectType {
	return envpkg.Flower
}
