package systems

import (
	"math"

	"github.com/deglan/horrorchain/engine/entities"
)

func CollectPotions(player *entities.Player, potions []*entities.Potion) []*entities.Potion {
	var remaining []*entities.Potion
	for _, potion := range potions {
		dx := player.X - potion.X
		dy := player.Y - potion.Y
		if math.Hypot(dx, dy) < 4.0 {
			player.CombatComp.HealthUpdate(potion.AmtHeal)
		} else {
			remaining = append(remaining, potion)
		}
	}
	return remaining
}
