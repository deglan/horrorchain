package entities

import (
	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/components"
)

type PlayerState uint8

const (
	Down PlayerState = iota
	Up
	Left
	Right
)

type Player struct {
	*Sprite
	Health     uint
	Animations map[PlayerState]*animations.Animation
	CombatComp *components.BasicCombat
}

func (p *Player) ActiveAnimation(dx, dy int) *animations.Animation {
	if dx == 0 && dy == 0 {
		return p.Animations[Down]
	}

	if dx < 0 {
		return p.Animations[Left]
	}

	if dx > 0 {
		return p.Animations[Right]
	}

	if dy < 0 {
		return p.Animations[Up]
	}

	return p.Animations[Down]
}
