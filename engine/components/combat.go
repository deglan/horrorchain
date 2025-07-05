package components

type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{health, attackPower, false}
}

func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

func (b *BasicCombat) Health() int {
	return b.health
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) Attack() bool {
	return true
}

func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Update() {
	b.attacking = false
}

var _ Combat = (*BasicCombat)(nil)

type EnemyCombat struct {
	*BasicCombat
	attackCooldown      int
	timeSinceLastAttack int
}

func NewEnemyCombat(health, attackPower, attackCooldown int) *EnemyCombat {
	return &EnemyCombat{NewBasicCombat(health, attackPower), attackCooldown, 0}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceLastAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceLastAttack = 0
		return true
	}
	return false
}

func (e *EnemyCombat) Update() {
	e.timeSinceLastAttack += 1
	e.attacking = false
}

var _ Combat = (*EnemyCombat)(nil)
