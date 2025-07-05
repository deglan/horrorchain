package systems

import (
	"fmt"

	"github.com/deglan/horrorchain/engine/entities"
)

func CollectPotions(player *entities.Player, potions []*entities.Potion) {
	for _, potion := range potions {
		if player.X > potion.X {
			player.Health += potion.AmtHeal
			fmt.Printf("Picked up potion! Health: %d\n", player.Health)
		}
	}
}
