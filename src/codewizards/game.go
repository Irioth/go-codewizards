package codewizards

type Game struct {
	RandomSeed                           int64
	TickCount                            int
	MapSize                              float64
	SkillsEnabled                        bool
	RawMessagesEnabled                   bool
	FriendlyFireDamageFactor             float64
	BuildingDamageScoreFactor            float64
	BuildingEliminationScoreFactor       float64
	MinionDamageScoreFactor              float64
	MinionEliminationScoreFactor         float64
	WizardDamageScoreFactor              float64
	WizardEliminationScoreFactor         float64
	TeamWorkingScoreFactor               float64
	VictoryScore                         int
	ScoreGainRange                       float64
	RawMessageMaxLength                  int
	RawMessageTransmissionSpeed          float64
	WizardRadius                         float64
	WizardCastRange                      float64
	WizardVisionRange                    float64
	WizardForwardSpeed                   float64
	WizardBackwardSpeed                  float64
	WizardStrafeSpeed                    float64
	WizardBaseLife                       int
	WizardLifeGrowthPerLevel             int
	WizardBaseMana                       int
	WizardManaGrowthPerLevel             int
	WizardBaseLifeRegeneration           float64
	WizardLifeRegenerationGrowthPerLevel float64
	WizardBaseManaRegeneration           float64
	WizardManaRegenerationGrowthPerLevel float64
	WizardMaxTurnAngle                   float64
	WizardMaxResurrectionDelayTicks      int
	WizardMinResurrectionDelayTicks      int
	WizardActionCooldownTicks            int
	StaffCooldownTicks                   int
	MagicMissileCooldownTicks            int
	FrostBoltCooldownTicks               int
	FireballCooldownTicks                int
	HasteCooldownTicks                   int
	ShieldCooldownTicks                  int
	MagicMissileManacost                 int
	FrostBoltManacost                    int
	FireballManacost                     int
	HasteManacost                        int
	ShieldManacost                       int
	StaffDamage                          int
	StaffSector                          float64
	StaffRange                           float64
	LevelUpXpValues                      []int
	MinionRadius                         float64
	MinionVisionRange                    float64
	MinionSpeed                          float64
	MinionMaxTurnAngle                   float64
	MinionLife                           int
	FactionMinionAppearanceIntervalTicks int
	OrcWoodcutterActionCooldownTicks     int
	OrcWoodcutterDamage                  int
	OrcWoodcutterAttackSector            float64
	OrcWoodcutterAttackRange             float64
	FetishBlowdartActionCooldownTicks    int
	FetishBlowdartAttackRange            float64
	FetishBlowdartAttackSector           float64
	BonusRadius                          float64
	BonusAppearanceIntervalTicks         int
	BonusScoreAmount                     int
	DartRadius                           float64
	DartSpeed                            float64
	DartDirectDamage                     int
	MagicMissileRadius                   float64
	MagicMissileSpeed                    float64
	MagicMissileDirectDamage             int
	FrostBoltRadius                      float64
	FrostBoltSpeed                       float64
	FrostBoltDirectDamage                int
	FireballRadius                       float64
	FireballSpeed                        float64
	FireballExplosionMaxDamageRange      float64
	FireballExplosionMinDamageRange      float64
	FireballExplosionMaxDamage           int
	FireballExplosionMinDamage           int
	GuardianTowerRadius                  float64
	GuardianTowerVisionRange             float64
	GuardianTowerLife                    float64
	GuardianTowerAttackRange             float64
	GuardianTowerDamage                  int
	GuardianTowerCooldownTicks           int
	FactionBaseRadius                    float64
	FactionBaseVisionRange               float64
	FactionBaseLife                      float64
	FactionBaseAttackRange               float64
	FactionBaseDamage                    int
	FactionBaseCooldownTicks             int
	BurningDurationTicks                 int
	BurningSummaryDamage                 int
	EmpoweredDurationTicks               int
	EmpoweredDamageFactor                float64
	FrozenDurationTicks                  int
	HastenedDurationTicks                int
	HastenedBonusDurationFactor          float64
	HastenedMovementBonusFactor          float64
	HastenedRotationBonusFactor          float64
	ShieldedDurationTicks                int
	ShieldedBonusDurationFactor          float64
	ShieldedDirectDamageAbsorptionFactor float64
	AuraSkillRange                       float64
	RangeBonusPerSkillLevel              float64
	MagicalDamageBonusPerSkillLevel      int
	StaffDamageBonusPerSkillLevel        int
	MovementBonusFactorPerSkillLevel     float64
	MagicalDamageAbsorptionPerSkillLevel int
}