package codewizards

type Strategy interface {
	Move(me *Wizard, world *World, game *Game, move *Move)
}

type Faction int

const (
	Faction_Academy Faction = iota
	Faction_Renegades
	Faction_Neutral
	Faction_Other
)

type Unit struct {
	Id             int64
	X, Y           float64
	SpeedX, SpeedY float64
	Angle          float64
	Faction        Faction
}

type CircularUnit struct {
	Unit
	Radius float64
}

type LivingUnit struct {
	CircularUnit
	Life, MaxLife int
	Statuses      []*Status
}

type StatusType int

const (
	Status_Burning StatusType = iota
	Status_Empowered
	Status_Frozen
	Status_Hastened
	Status_Shielded
)

type Status struct {
	Id                     int64
	Type                   StatusType
	WizardId               int64
	PlayerId               int64
	RemainingDurationTicks int
}
