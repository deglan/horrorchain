package entities

import "github.com/deglan/horrorchain/engine/components"

type Enemy struct {
	*Sprite
	FollowsPlayer bool
	CombatComp    *components.EnemyCombat
}
