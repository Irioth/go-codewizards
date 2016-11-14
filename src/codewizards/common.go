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

// Unit -----------------------------------------------------------------------
type Unit struct {
	Id             int64
	X, Y           float64
	SpeedX, SpeedY float64
	Angle          float64
	Faction        Faction
}

func (u *Unit) GetId() int64        { return u.Id }
func (u *Unit) GetX() float64       { return u.X }
func (u *Unit) GetY() float64       { return u.Y }
func (u *Unit) GetSpeedX() float64  { return u.SpeedX }
func (u *Unit) GetSpeedY() float64  { return u.SpeedY }
func (u *Unit) GetAngle() float64   { return u.Angle }
func (u *Unit) GetFaction() Faction { return u.Faction }
func (u *Unit) AsUnit() *Unit       { return u }

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

type Point interface {
	GetX() float64
	GetY() float64
}

func (u *Unit) GetAngleToPoint(p Point) float64 {
	return u.GetAngleTo(p.GetX(), p.GetY())
}

func (u *Unit) GetDistanceToPoint(p Point) float64 {
	return u.GetDistanceTo(p.GetX(), p.GetY())
}

// CircularUnit ---------------------------------------------------------------
type CircularUnit struct {
	Unit
	Radius float64
}

func (u *CircularUnit) GetRadius() float64            { return u.Radius }
func (u *CircularUnit) AsCircularUnit() *CircularUnit { return u }

// LivingUnit -----------------------------------------------------------------
type LivingUnit struct {
	CircularUnit
	Life, MaxLife int
	Statuses      []*Status
}

func (u *LivingUnit) GetLife() int              { return u.Life }
func (u *LivingUnit) GetMaxLife() int           { return u.MaxLife }
func (u *LivingUnit) GetStatuses() []*Status    { return u.Statuses }
func (u *LivingUnit) AsLivingUnit() *LivingUnit { return u }

// Statuses -------------------------------------------------------------------
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
