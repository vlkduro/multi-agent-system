package object

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

type Hive struct {
	id       envpkg.ObjectID
	pos      *envpkg.Position
	qHoney   int
	qNectar  int
	qPollen  int
	queen    bool
	minHoney int
}

type HiveJson struct {
	ID             string          `json:"id"`
	Position       envpkg.Position `json:"position"`
	QuantityHoney  int             `json:"quantity_honey"`
	QuantityNectar int             `json:"quantity_nectar"`
	QuantityPollen int             `json:"quantity_pollen"`
	Quenn          bool            `json:"queen"`
	MinHoney       int             `json:"min_honey"`
}

func NewHive(id string, pos *envpkg.Position, qHoney int, qNectar int, qPollen int, minHoney int) *Hive {
	return &Hive{
		id:       envpkg.ObjectID(id),
		pos:      pos.Copy(),
		qHoney:   qHoney,
		qNectar:  qNectar,
		qPollen:  qPollen,
		queen:    true,
		minHoney: minHoney,
	}
}

func (h Hive) ID() envpkg.ObjectID {
	return h.id
}

func (h Hive) Position() *envpkg.Position {
	return h.pos.Copy()
}

func (h Hive) Copy() interface{} {
	return &Hive{
		id:       h.id,
		pos:      h.pos.Copy(),
		qHoney:   h.qHoney,
		qNectar:  h.qNectar,
		qPollen:  h.qPollen,
		queen:    h.queen,
		minHoney: h.minHoney,
	}
}

/*func (h Hive) Become(interface{}) {

}

func (h Hive) ToJsonObj() interface{}
func (h Hive) Interact()

func (h Hive) checkResourcesAndBees() {

}

func (h Hive) growColony() {

}

func (h Hive) Die() {

}*/
