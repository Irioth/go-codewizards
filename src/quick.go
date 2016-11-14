package main

import (
	. "codewizards"
	"codewizards/runner"
	"math"
	"math/rand"
)

func main() {
	runner.Start(NewQuick)
}

const (
	WayPointsRadius = 100.0
	LowHPFactor     = 0.25
)

var (
	WayPoints = map[LaneType][]point{
		Lane_Middle: []point{
			{100, 4000 - 100},
			{600, 4000 - 200},
			// {200, 4000 - 600),
			{800, 4000 - 800},
			{4000 - 600, 600},
		},
		Lane_Top: []point{
			{100, 4000 - 100},
			{100, 4000 - 400},
			{200, 4000 - 800},
			{200, 4000 * 0.75},
			{200, 4000 * 0.5},
			{200, 4000 * 0.25},
			{200, 200},
			{4000 * 0.25, 200},
			{4000 * 0.5, 200},
			{4000 * 0.75, 200},
			{4000 - 200, 200},
		},
		Lane_Bottom: []point{
			{100, 4000 - 100},
			{400, 4000 - 100},
			{800, 4000 - 200},
			{4000 * 0.25, 4000 - 200},
			{4000 * 0.5, 4000 - 200},
			{4000 * 0.75, 4000 - 200},
			{4000 - 200, 4000 - 200},
			{4000 - 200, 4000 * 0.75},
			{4000 - 200, 4000 * 0.5},
			{4000 - 200, 4000 * 0.25},
			{4000 - 200, 200},
		},
	}
)

type Context struct {
	me    *Wizard
	world *World
	game  *Game
	move  *Move
}

func (c *Context) initializeTick(me *Wizard, world *World, game *Game, move *Move) {
	c.me = me
	c.world = world
	c.game = game
	c.move = move
}

type QuickStrategy struct {
	Context
	random      *rand.Rand
	waypoints   []point
	rewaypoints []point
}

func NewQuick() Strategy {
	return &QuickStrategy{}
}

func (s *QuickStrategy) Move(me *Wizard, world *World, game *Game, move *Move) {
	s.initializeStrategy(me, game)
	s.initializeTick(me, world, game, move)

	move.StrafeSpeed = game.WizardStrafeSpeed
	if s.random.Intn(2) == 0 {
		move.StrafeSpeed = -move.StrafeSpeed
	}

	if float64(me.Life) < LowHPFactor*float64(me.MaxLife) {
		s.goTo(s.previousWaypoint())
		return
	}

	nearestTarget := getNearestTarget(me, world)

	if nearestTarget != nil {
		distance := me.GetDistanceToUnit(nearestTarget.AsUnit())

		if distance <= me.CastRange {
			angle := me.GetAngleToUnit(nearestTarget.AsUnit())

			move.Turn = angle

			if math.Abs(angle) < game.StaffSector/2. {
				move.Action = Action_MagicMissle
				move.CastAngle = angle
				move.MinCastDistance = distance - nearestTarget.Radius + game.MagicMissileRadius
			}

			return
		}
	}

	s.goTo(s.nextWaypoint())
}

func (s *QuickStrategy) initializeStrategy(me *Wizard, game *Game) {
	if s.random == nil {
		s.random = rand.New(rand.NewSource(game.RandomSeed))
		switch me.Id {
		case 1, 2, 6, 7:
			s.waypoints = WayPoints[Lane_Top]
		case 3, 8:
			s.waypoints = WayPoints[Lane_Middle]
		case 4, 5, 9, 10:
			s.waypoints = WayPoints[Lane_Bottom]
		}

		// calc reversed version on waypoints
		l := len(s.waypoints)
		s.rewaypoints = make([]point, l)
		for i, p := range s.waypoints {
			s.rewaypoints[l-1-i] = p
		}

	}
}

func getNearestTarget(me *Wizard, world *World) *LivingUnit {
	targets := make([]*LivingUnit, 0)

	for _, b := range world.Buildings {
		targets = append(targets, &b.LivingUnit)
	}
	for _, w := range world.Wizards {
		targets = append(targets, &w.LivingUnit)
	}
	for _, m := range world.Minions {
		targets = append(targets, &m.LivingUnit)
	}

	var nearestTarget *LivingUnit = nil
	nearestTargetDistance := math.MaxFloat64

	for _, t := range targets {
		if t.Faction == Faction_Neutral || t.Faction == me.Faction {
			continue
		}

		distance := me.GetDistanceToUnit(t.AsUnit())

		if distance < nearestTargetDistance {
			nearestTarget = t
			nearestTargetDistance = distance
		}
	}

	return nearestTarget
}

func (s *QuickStrategy) nextWaypoint() point {
	return s.findNextWaypoint(s.waypoints)
}

func (s *QuickStrategy) previousWaypoint() point {
	return s.findNextWaypoint(s.rewaypoints)
}

func (s *QuickStrategy) findNextWaypoint(waypoints []point) point {
	last := waypoints[len(waypoints)-1]
	for i, p := range waypoints[:len(waypoints)-1] {
		if p.Distance(s.me.X, s.me.Y) <= 100 {
			return waypoints[i+1]
		}

		if last.Distance(p.x, p.y) < last.Distance(s.me.X, s.me.Y) {
			return p
		}
	}
	return last
}

func (s *QuickStrategy) goTo(p point) {
	angle := s.me.GetAngleTo(p.x, p.y)

	s.move.Turn = angle

	if math.Abs(angle) < s.game.StaffSector/4 {
		s.move.Speed = s.game.WizardForwardSpeed
	}

}

type point struct {
	x, y float64
}

func (p point) Distance(x, y float64) float64 {
	return math.Hypot(p.x-x, p.y-y)
}
