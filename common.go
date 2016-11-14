package codewizards

import (
	"math"
)

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

func (u *Unit) GetAngleTo(x, y float64) float64 {
	absolute_angle_to := math.Atan2(y-u.Y, x-u.X)
	relative_angle_to := absolute_angle_to - u.Angle

	for relative_angle_to > math.Pi {
		relative_angle_to -= 2.0 * math.Pi
	}
	for relative_angle_to < -math.Pi {
		relative_angle_to += 2.0 * math.Pi
	}
	return relative_angle_to
}

func (u *Unit) GetAngleToUnit(unit *Unit) float64 {
	return u.GetAngleTo(unit.X, unit.Y)
}

func (u *Unit) GetDistanceTo(x, y float64) float64 {
	return math.Hypot(x-u.X, y-u.Y)
}

func (u *Unit) GetDistanceToUnit(unit *Unit) float64 {
	return u.GetDistanceTo(unit.X, unit.Y)
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
