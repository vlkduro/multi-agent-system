package object

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

type Hive struct {
	id          envpkg.ObjectID
	Pos         *envpkg.Position
	qHoney      int
	qNectar     int
	qPollen     int
	queen       bool
	minHoney    int
	env         *envpkg.Environment
	flowerStack *utils.Stack[*Flower]
}

type HiveJson struct {
	ID             string          `json:"id"`
	Position       envpkg.Position `json:"position"`
	QuantityHoney  int             `json:"quantity_honey"`
	QuantityNectar int             `json:"quantity_nectar"`
	QuantityPollen int             `json:"quantity_pollen"`
	Queen          bool            `json:"queen"`
	MinHoney       int             `json:"min_honey"`
}

func NewHive(id string, pos *envpkg.Position, qHoney int, qNectar int, qPollen int, minHoney int, environment *envpkg.Environment) *Hive {
	return &Hive{
		id:          envpkg.ObjectID(id),
		Pos:         pos.Copy(),
		qHoney:      qHoney,
		qNectar:     qNectar,
		qPollen:     qPollen,
		queen:       true,
		minHoney:    minHoney,
		env:         environment,
		flowerStack: utils.NewStack[*Flower](),
	}
}

func (h Hive) ID() envpkg.ObjectID {
	return h.id
}

func (h Hive) Position() *envpkg.Position {
	return h.Pos.Copy()
}

func (h Hive) Copy() interface{} {
	return &Hive{
		id:       h.id,
		Pos:      h.Pos.Copy(),
		qHoney:   h.qHoney,
		qNectar:  h.qNectar,
		qPollen:  h.qPollen,
		queen:    h.queen,
		minHoney: h.minHoney,
	}
}

func (h *Hive) Become(h_alt interface{}) {
	if h_alt == nil {
		return
	}

	altered_hive, ok := h_alt.(*Hive)

	if ok {
		h.id = altered_hive.id
		h.Pos = altered_hive.Pos.Copy()
		h.qHoney = altered_hive.qHoney
		h.qNectar = altered_hive.qNectar
		h.qPollen = altered_hive.qPollen
		h.queen = altered_hive.queen
		h.minHoney = altered_hive.minHoney
	}
}

func (h *Hive) Update() {
	nb_bees := h.env.GetNumberAgent() - utils.GetNumberHornets()
	if h.qNectar > nb_bees {
		h.qNectar -= nb_bees / 2
		h.qNectar -= 4
		h.qHoney += 1
	} /*else {
		h.Die()
	}*/
	// Compter les abeilles de la ruche pour débiter du miel
	// Ou si autre mieux, le faire
}

func (h Hive) ToJsonObj() interface{} {
	return HiveJson{
		ID:             string(h.id),
		Position:       *h.Pos.Copy(),
		QuantityHoney:  h.qHoney,
		QuantityNectar: h.qNectar,
		QuantityPollen: h.qPollen,
		Queen:          h.queen,
		MinHoney:       h.minHoney,
	}
}

func (h Hive) GetPosition() (position envpkg.Position) {
	position = *h.Pos.Copy()
	return
}

func (h *Hive) StoreNectar(nectar int) {
	h.qNectar += nectar
}

func (h *Hive) GetNectar(nectar int) {
	h.qNectar -= nectar
}

func (h *Hive) StoreHoney(honey int) {
	h.qHoney += honey
}

func (h Hive) GetQNectar() int {
	return h.qNectar
}

func (h Hive) IsAlive() (isAlive bool) {
	if h.qHoney >= h.minHoney && h.queen {
		isAlive = true
	}
	return
}

/*func (h Hive) Type() envpkg.AgentType {
	return envpkg.Hive
}*/

func (h *Hive) Die() {
	h.queen = false
	h.env.RemoveObject(h)
	h.IsAlive()
	h.Pos = nil
}

func (h Hive) TypeObject() envpkg.ObjectType {
	return envpkg.Hive
}

func (h Hive) GetHoney() int {
	return h.qHoney
}

func (h *Hive) RetreiveHoney(honey int) {
	h.qHoney -= honey
}
func (h *Hive) AddFlower(flower *Flower) {
	h.flowerStack.Push(flower)
}

func (h *Hive) GetLatestFlowerPos() *Flower {
	if h.flowerStack.IsEmpty() {
		return nil
	}
	flower, _ := h.flowerStack.Pop()
	return flower
}
