package codewizards

type ActionType int

const (
	Action_None ActionType = iota
	Action_Staff
	Action_MagicMissle
	Action_FrostBolt
	Action_Fireball
	Action_Haste
	Action_Shield
)

type Move struct {
	Speed, StrafeSpeed float64
	Turn               float64
	Action             ActionType
	CastAngle          float64
	MinCastDistance    float64
	MaxCastDistance    float64
	StatusTargetId     int64
	SkillToLearn       SkillType
	Messages           []*Message
}

func NewMove() *Move {
	return &Move{
		MaxCastDistance: 10000,
		StatusTargetId:  -1,
	}
}
