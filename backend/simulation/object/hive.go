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

func (h Hive) Become(h_alt interface{}) {
	if h_alt == nil {
		return
	}

	altered_hive, ok := h_alt.(*Hive)

	if ok {
		h.id = altered_hive.id
		h.pos = altered_hive.pos.Copy()
		h.qHoney = altered_hive.qHoney
		h.qNectar = altered_hive.qNectar
		h.qPollen = altered_hive.qPollen
		h.queen = altered_hive.queen
		h.minHoney = altered_hive.minHoney
	}
}

func (h Hive) ToJsonObj() interface{} {
	return HiveJson{
		ID:             string(h.id),
		Position:       *h.pos.Copy(),
		QuantityHoney:  h.qHoney,
		QuantityNectar: h.qNectar,
		QuantityPollen: h.qPollen,
		Quenn:          h.queen,
		MinHoney:       h.minHoney,
	}
}

func (h Hive) StoreNectar(nectar int) {
	h.qNectar += nectar
}

func (h Hive) IsAlive() (isAlive bool) {
	if h.qHoney >= h.minHoney && h.queen {
		isAlive = true
	}
	return
}

func (h *Hive) Die() {
	h.queen = false
}
