package codewizards

type World struct {
	TickIndex     int
	TickCount     int
	Width, Height float64
	Players       []*Player
	Wizards       []*Wizard
	Minions       []*Minion
	Projectiles   []*Projectile
	Bonuses       []*Bonus
	Buildings     []*Building
	Trees         []*Tree
}

type Tree struct {
	LivingUnit
}

// Projectiles ----------------------------------------------------------------
type ProjectileType int

const (
	Projectile_MagicMissle ProjectileType = iota
	Projectile_FrostBolt
	Projectile_Fireball
	Projectile_Dart
)

type Projectile struct {
	CircularUnit
	Type          ProjectileType
	OwnerUnitId   int64
	OwnerPlayerId int64
}

// Bonuses --------------------------------------------------------------------
type BonusType int

const (
	Bonus_Empower BonusType = iota
	Bonus_Haste
	Bonus_Shield
)

type Bonus struct {
	CircularUnit
	Type BonusType
}

// Buildings ------------------------------------------------------------------
type BuildingType int

const (
	Building_GuardianTower BuildingType = iota
	Building_FactionBase
)

type Building struct {
	LivingUnit
	Type                         BuildingType
	VisionRange                  float64
	AttackRange                  float64
	Damage                       int
	CooldownTicks                int
	RemainingActionCooldownTicks int
}

// Minions --------------------------------------------------------------------
type MinionType int

const (
	Minion_OrcWoodcutter MinionType = iota
	Minion_FetishBlowdart
)

type Minion struct {
	LivingUnit
	Type                         MinionType
	VisionRange                  float64
	Damage                       int
	CooldownTicks                int
	RemainingActionCooldownTicks int
}
