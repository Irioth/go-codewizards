package runner

import (
	"bufio"
	"encoding/binary"
	"errors"
	. "codewizards"
	"net"
)

var (
	Order = binary.LittleEndian
)

type PlayerContext struct {
	Wizards []*Wizard
	World   *World
}

type Client struct {
	conn          net.Conn
	w             *bufio.Writer
	r             *bufio.Reader
	previousTrees []*Tree
}

type MessageType int

const (
	Message_Unknown MessageType = iota
	Message_GameOver
	Message_AuthToken
	Message_TeamSize
	Message_ProtoVersion
	Message_GameContext
	Message_PlayerContext
	Message_Moves
)

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn, bufio.NewWriter(conn), bufio.NewReader(conn), nil}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) WriteToken(token string) {
	c.writeOpcode(Message_AuthToken)
	c.writeString(token)
	c.flush()
}

func (c *Client) WriteProtocolVersion() {
	c.writeOpcode(Message_ProtoVersion)
	c.writeInt(1)
	c.flush()
}

func (c *Client) WriteMoves(moves []*Move) {
	c.writeOpcode(Message_Moves)
	c.writeInt(len(moves))
	for i := range moves {
		c.writeMove(moves[i])
	}
	c.flush()
}

func (c *Client) writeMove(move *Move) {
	c.writeByte(1)
	c.writeFloat64(move.Speed)
	c.writeFloat64(move.StrafeSpeed)
	c.writeFloat64(move.Turn)
	c.writeByte(byte(move.Action))
	c.writeFloat64(move.CastAngle)
	c.writeFloat64(move.MinCastDistance)
	c.writeFloat64(move.MaxCastDistance)
	c.writeInt64(move.StatusTargetId)
	c.writeByte(byte(move.SkillToLearn))
	c.writeInt(len(move.Messages))
	for _, m := range move.Messages {
		c.writeByte(1)
		c.writeByte(byte(m.Lane))
		c.writeByte(byte(m.SkillToLearn))
		c.writeBytes(m.RawMessage)
	}
}

func (c *Client) ReadTeamSize() int {
	c.ensureMessageType(c.readByte(), Message_TeamSize)
	return c.readInt()
}

func (c *Client) ReadPlayerContext() *PlayerContext {
	opcode := c.readByte()
	if opcode == byte(Message_GameOver) {
		return nil
	}
	c.ensureMessageType(opcode, Message_PlayerContext)
	if !c.readBool() {
		return nil
	}
	return &PlayerContext{
		Wizards: c.readWizards(),
		World:   c.readWorld(),
	}
}

func (c *Client) readWorld() *World {
	if !c.readBool() {
		return nil
	}
	return &World{
		TickIndex:   c.readInt(),
		TickCount:   c.readInt(),
		Width:       c.readFloat64(),
		Height:      c.readFloat64(),
		Players:     c.readPlayers(),
		Wizards:     c.readWizards(),
		Minions:     c.readMinions(),
		Projectiles: c.readProjectiles(),
		Bonuses:     c.readBonuses(),
		Buildings:   c.readBuildings(),
		Trees:       c.readTrees(),
	}
}

func (c *Client) readPlayers() []*Player {
	l := c.readInt()
	r := make([]*Player, l)
	for i := range r {
		r[i] = c.readPlayer()
	}
	return r
}

func (c *Client) readPlayer() *Player {
	if !c.readBool() {
		return nil
	}
	return &Player{
		Id:              c.readInt64(),
		Me:              c.readBool(),
		Name:            c.readString(),
		StrategyCrashed: c.readBool(),
		Score:           c.readInt(),
		Faction:         Faction(c.readByte()),
	}
}

func (c *Client) readMinions() []*Minion {
	l := c.readInt()
	r := make([]*Minion, l)
	for i := range r {
		r[i] = c.readMinion()
	}
	return r
}

func (c *Client) readMinion() *Minion {
	if !c.readBool() {
		return nil
	}
	return &Minion{
		LivingUnit:                   c.readLivingUnit(),
		Type:                         MinionType(c.readByte()),
		VisionRange:                  c.readFloat64(),
		Damage:                       c.readInt(),
		CooldownTicks:                c.readInt(),
		RemainingActionCooldownTicks: c.readInt(),
	}
}

func (c *Client) readBuildings() []*Building {
	l := c.readInt()
	r := make([]*Building, l)
	for i := range r {
		r[i] = c.readBuilding()
	}
	return r
}

func (c *Client) readBuilding() *Building {
	if !c.readBool() {
		return nil
	}
	return &Building{
		LivingUnit:                   c.readLivingUnit(),
		Type:                         BuildingType(c.readByte()),
		VisionRange:                  c.readFloat64(),
		AttackRange:                  c.readFloat64(),
		Damage:                       c.readInt(),
		CooldownTicks:                c.readInt(),
		RemainingActionCooldownTicks: c.readInt(),
	}
}

func (c *Client) readTrees() []*Tree {
	l := c.readInt()
	if l < 0 {
		return c.previousTrees
	}
	r := make([]*Tree, l)
	for i := range r {
		r[i] = c.readTree()
	}
	c.previousTrees = r
	return r
}

func (c *Client) readTree() *Tree {
	if !c.readBool() {
		return nil
	}
	return &Tree{
		LivingUnit: c.readLivingUnit(),
	}
}

func (c *Client) readProjectiles() []*Projectile {
	l := c.readInt()
	r := make([]*Projectile, l)
	for i := range r {
		r[i] = c.readProjectile()
	}
	return r
}

func (c *Client) readProjectile() *Projectile {
	if !c.readBool() {
		return nil
	}
	return &Projectile{
		CircularUnit:  c.readCircularUnit(),
		Type:          ProjectileType(c.readByte()),
		OwnerUnitId:   c.readInt64(),
		OwnerPlayerId: c.readInt64(),
	}
}

func (c *Client) readBonuses() []*Bonus {
	l := c.readInt()
	r := make([]*Bonus, l)
	for i := range r {
		r[i] = c.readBonus()
	}
	return r
}

func (c *Client) readBonus() *Bonus {
	if !c.readBool() {
		return nil
	}
	return &Bonus{
		CircularUnit: c.readCircularUnit(),
		Type:         BonusType(c.readByte()),
	}
}

func (c *Client) readWizards() []*Wizard {
	l := c.readInt()
	r := make([]*Wizard, l)
	for i := range r {
		r[i] = c.readWizard()
	}
	return r
}

func (c *Client) readWizard() *Wizard {
	if !c.readBool() {
		return nil
	}
	return &Wizard{
		LivingUnit:    c.readLivingUnit(),
		OwnerPlayerId: c.readInt64(),
		Me:            c.readBool(),
		Mana:          c.readInt(),
		MaxMana:       c.readInt(),
		VisionRange:   c.readFloat64(),
		CastRange:     c.readFloat64(),
		Xp:            c.readInt(),
		Level:         c.readInt(),
		Skills:        c.readSkillTypes(),
		RemainingActionCooldownTicks:   c.readInt(),
		RemainingCooldownTicksByAction: c.readIntArray(),
		Master:   c.readBool(),
		Messages: c.readMessages(),
	}
}

func (c *Client) readMessages() []*Message {
	l := c.readInt()
	r := make([]*Message, l)
	for i := range r {
		r[i] = c.readMessage()
	}
	return r
}

func (c *Client) readMessage() *Message {
	if !c.readBool() {
		return nil
	}
	return &Message{
		Lane:         LaneType(c.readByte()),
		SkillToLearn: SkillType(c.readByte()),
		RawMessage:   c.readBytes(),
	}
}

func (c *Client) readBytes() []byte {
	l := c.readInt()
	r := make([]byte, l)
	for i := range r {
		r[i] = c.readByte()
	}
	return r
}

func (c *Client) readSkillTypes() []SkillType {
	l := c.readInt()
	r := make([]SkillType, l)
	for i := range r {
		r[i] = SkillType(c.readByte())
	}
	return r
}

func (c *Client) readLivingUnit() LivingUnit {
	return LivingUnit{
		CircularUnit: c.readCircularUnit(),
		Life:         c.readInt(),
		MaxLife:      c.readInt(),
		Statuses:     c.readStatuses(),
	}
}
func (c *Client) readCircularUnit() CircularUnit {
	return CircularUnit{
		Unit: Unit{
			Id:      c.readInt64(),
			X:       c.readFloat64(),
			Y:       c.readFloat64(),
			SpeedX:  c.readFloat64(),
			SpeedY:  c.readFloat64(),
			Angle:   c.readFloat64(),
			Faction: Faction(c.readByte()),
		},
		Radius: c.readFloat64(),
	}
}

func (c *Client) readStatuses() []*Status {
	l := c.readInt()
	r := make([]*Status, l)
	for i := range r {
		r[i] = c.readStatus()
	}
	return r
}

func (c *Client) readStatus() *Status {
	return &Status{
		Id:                     c.readInt64(),
		Type:                   StatusType(c.readByte()),
		WizardId:               c.readInt64(),
		PlayerId:               c.readInt64(),
		RemainingDurationTicks: c.readInt(),
	}

}

func (c *Client) ReadGameContext() *Game {
	c.ensureMessageType(c.readByte(), Message_GameContext)
	if !c.readBool() {
		return nil
	}
	return &Game{
		RandomSeed:                           c.readInt64(),
		TickCount:                            c.readInt(),
		MapSize:                              c.readFloat64(),
		SkillsEnabled:                        c.readBool(),
		RawMessagesEnabled:                   c.readBool(),
		FriendlyFireDamageFactor:             c.readFloat64(),
		BuildingDamageScoreFactor:            c.readFloat64(),
		BuildingEliminationScoreFactor:       c.readFloat64(),
		MinionDamageScoreFactor:              c.readFloat64(),
		MinionEliminationScoreFactor:         c.readFloat64(),
		WizardDamageScoreFactor:              c.readFloat64(),
		WizardEliminationScoreFactor:         c.readFloat64(),
		TeamWorkingScoreFactor:               c.readFloat64(),
		VictoryScore:                         c.readInt(),
		ScoreGainRange:                       c.readFloat64(),
		RawMessageMaxLength:                  c.readInt(),
		RawMessageTransmissionSpeed:          c.readFloat64(),
		WizardRadius:                         c.readFloat64(),
		WizardCastRange:                      c.readFloat64(),
		WizardVisionRange:                    c.readFloat64(),
		WizardForwardSpeed:                   c.readFloat64(),
		WizardBackwardSpeed:                  c.readFloat64(),
		WizardStrafeSpeed:                    c.readFloat64(),
		WizardBaseLife:                       c.readInt(),
		WizardLifeGrowthPerLevel:             c.readInt(),
		WizardBaseMana:                       c.readInt(),
		WizardManaGrowthPerLevel:             c.readInt(),
		WizardBaseLifeRegeneration:           c.readFloat64(),
		WizardLifeRegenerationGrowthPerLevel: c.readFloat64(),
		WizardBaseManaRegeneration:           c.readFloat64(),
		WizardManaRegenerationGrowthPerLevel: c.readFloat64(),
		WizardMaxTurnAngle:                   c.readFloat64(),
		WizardMaxResurrectionDelayTicks:      c.readInt(),
		WizardMinResurrectionDelayTicks:      c.readInt(),
		WizardActionCooldownTicks:            c.readInt(),
		StaffCooldownTicks:                   c.readInt(),
		MagicMissileCooldownTicks:            c.readInt(),
		FrostBoltCooldownTicks:               c.readInt(),
		FireballCooldownTicks:                c.readInt(),
		HasteCooldownTicks:                   c.readInt(),
		ShieldCooldownTicks:                  c.readInt(),
		MagicMissileManacost:                 c.readInt(),
		FrostBoltManacost:                    c.readInt(),
		FireballManacost:                     c.readInt(),
		HasteManacost:                        c.readInt(),
		ShieldManacost:                       c.readInt(),
		StaffDamage:                          c.readInt(),
		StaffSector:                          c.readFloat64(),
		StaffRange:                           c.readFloat64(),
		LevelUpXpValues:                      c.readIntArray(),
		MinionRadius:                         c.readFloat64(),
		MinionVisionRange:                    c.readFloat64(),
		MinionSpeed:                          c.readFloat64(),
		MinionMaxTurnAngle:                   c.readFloat64(),
		MinionLife:                           c.readInt(),
		FactionMinionAppearanceIntervalTicks: c.readInt(),
		OrcWoodcutterActionCooldownTicks:     c.readInt(),
		OrcWoodcutterDamage:                  c.readInt(),
		OrcWoodcutterAttackSector:            c.readFloat64(),
		OrcWoodcutterAttackRange:             c.readFloat64(),
		FetishBlowdartActionCooldownTicks:    c.readInt(),
		FetishBlowdartAttackRange:            c.readFloat64(),
		FetishBlowdartAttackSector:           c.readFloat64(),
		BonusRadius:                          c.readFloat64(),
		BonusAppearanceIntervalTicks:         c.readInt(),
		BonusScoreAmount:                     c.readInt(),
		DartRadius:                           c.readFloat64(),
		DartSpeed:                            c.readFloat64(),
		DartDirectDamage:                     c.readInt(),
		MagicMissileRadius:                   c.readFloat64(),
		MagicMissileSpeed:                    c.readFloat64(),
		MagicMissileDirectDamage:             c.readInt(),
		FrostBoltRadius:                      c.readFloat64(),
		FrostBoltSpeed:                       c.readFloat64(),
		FrostBoltDirectDamage:                c.readInt(),
		FireballRadius:                       c.readFloat64(),
		FireballSpeed:                        c.readFloat64(),
		FireballExplosionMaxDamageRange:      c.readFloat64(),
		FireballExplosionMinDamageRange:      c.readFloat64(),
		FireballExplosionMaxDamage:           c.readInt(),
		FireballExplosionMinDamage:           c.readInt(),
		GuardianTowerRadius:                  c.readFloat64(),
		GuardianTowerVisionRange:             c.readFloat64(),
		GuardianTowerLife:                    c.readFloat64(),
		GuardianTowerAttackRange:             c.readFloat64(),
		GuardianTowerDamage:                  c.readInt(),
		GuardianTowerCooldownTicks:           c.readInt(),
		FactionBaseRadius:                    c.readFloat64(),
		FactionBaseVisionRange:               c.readFloat64(),
		FactionBaseLife:                      c.readFloat64(),
		FactionBaseAttackRange:               c.readFloat64(),
		FactionBaseDamage:                    c.readInt(),
		FactionBaseCooldownTicks:             c.readInt(),
		BurningDurationTicks:                 c.readInt(),
		BurningSummaryDamage:                 c.readInt(),
		EmpoweredDurationTicks:               c.readInt(),
		EmpoweredDamageFactor:                c.readFloat64(),
		FrozenDurationTicks:                  c.readInt(),
		HastenedDurationTicks:                c.readInt(),
		HastenedBonusDurationFactor:          c.readFloat64(),
		HastenedMovementBonusFactor:          c.readFloat64(),
		HastenedRotationBonusFactor:          c.readFloat64(),
		ShieldedDurationTicks:                c.readInt(),
		ShieldedBonusDurationFactor:          c.readFloat64(),
		ShieldedDirectDamageAbsorptionFactor: c.readFloat64(),
		AuraSkillRange:                       c.readFloat64(),
		RangeBonusPerSkillLevel:              c.readFloat64(),
		MagicalDamageBonusPerSkillLevel:      c.readInt(),
		StaffDamageBonusPerSkillLevel:        c.readInt(),
		MovementBonusFactorPerSkillLevel:     c.readFloat64(),
		MagicalDamageAbsorptionPerSkillLevel: c.readInt(),
	}
}

func (c *Client) readIntArray() []int {
	count := c.readInt()
	r := make([]int, count)
	for i := range r {
		r[i] = c.readInt()
	}
	return r
}

func (c *Client) readInt() int {
	var v int32
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return int(v)
}

func (c *Client) readInt64() int64 {
	var v int64
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return v
}

func (c *Client) readFloat64() float64 {
	var v float64
	if err := binary.Read(c.r, Order, &v); err != nil {
		panic(err)
	}
	return v
}

func (c *Client) readBool() bool {
	return c.readByte() != 0
}

func (c *Client) readByte() byte {
	b, err := c.r.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (c *Client) readString() string {
	return string(c.readBytes())
}

func (c *Client) ensureMessageType(v byte, m MessageType) {
	if v != byte(m) {
		panic(errors.New("unexpected message"))
	}
}

func (c *Client) writeOpcode(m MessageType) {
	c.writeByte(byte(m))
}

func (c *Client) writeInt(v int) {
	if err := binary.Write(c.w, Order, int32(v)); err != nil {
		panic(err)
	}
}

func (c *Client) writeFloat64(v float64) {
	if err := binary.Write(c.w, Order, v); err != nil {
		panic(err)
	}
}

func (c *Client) writeInt64(v int64) {
	if err := binary.Write(c.w, Order, v); err != nil {
		panic(err)
	}
}

func (c *Client) writeByte(v byte) {
	if err := c.w.WriteByte(v); err != nil {
		panic(err)
	}
}

func (c *Client) writeBytes(v []byte) {
	c.writeInt(len(v))
	if _, err := c.w.Write(v); err != nil {
		panic(err)
	}
}

func (c *Client) writeString(v string) {
	c.writeInt(len(v))
	if _, err := c.w.WriteString(v); err != nil {
		panic(err)
	}
}

func (c *Client) flush() {
	if err := c.w.Flush(); err != nil {
		panic(err)
	}
}
