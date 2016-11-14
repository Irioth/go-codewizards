package codewizards

type Player struct {
	Id              int64
	Me              bool
	Name            string
	StrategyCrashed bool
	Score           int
	Faction         Faction
}

type SkillType int

const (
	Skill_RangeBonusPassive1 SkillType = iota
	Skill_RangeBonusAura1
	Skill_RangeBonusPassive2
	Skill_RangeBonusAura2
	Skill_AdvancedMagicMissle
	Skill_MagicalDamageBonusPassive1
	Skill_MagicalDamageBonusAura1
	Skill_MagicalDamageBonusPassive2
	Skill_MagicalDamageBonusAura2
	Skill_FrostBolt
	Skill_StaffDamageBonusPassive1
	Skill_StaffDamageBonusAura1
	Skill_StaffDamageBonusPassive2
	Skill_StaffDamageBonusAura2
	Skill_Fireball
	Skill_MovementBonusFactorPassive1
	Skill_MovementBonusFactorAura1
	Skill_MovementBonusFactorPassive2
	Skill_MovementBonusFactorAura2
	Skill_Haste
	Skill_MagicalDamageAbsorptionPassive1
	Skill_MagicalDamageAbsorptionAura1
	Skill_MagicalDamageAbsorptionPassive2
	Skill_MagicalDamageAbsorptionAura2
	Skill_Shield
)

type Wizard struct {
	LivingUnit
	OwnerPlayerId                  int64
	Me                             bool
	Mana, MaxMana                  int
	VisionRange                    float64
	CastRange                      float64
	Xp                             int
	Level                          int
	Skills                         []SkillType
	RemainingActionCooldownTicks   int
	RemainingCooldownTicksByAction []int
	Master                         bool
	Messages                       []*Message
}

type Message struct {
	Lane         LaneType
	SkillToLearn SkillType
	RawMessage   []byte
}

type LaneType int

const (
	Lane_Top LaneType = iota
	Lane_Middle
	Lane_Bottom
)
