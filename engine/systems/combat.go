package systems

import (
	"image"
	"math"

	"github.com/deglan/horrorchain/constants"
	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/entities"
)

func HandleCombat(player *entities.Player, enemies []*entities.Enemy, attackAnimations []*animations.AttackAnimation, colliders []image.Rectangle) []*entities.Enemy {
	deadEnemies := map[int]struct{}{}
	playerRect := image.Rect(
		int(player.X), int(player.Y),
		int(player.X)+constants.Tilesize, int(player.Y)+constants.Tilesize,
	)

	for idx, enemy := range enemies {
		enemy.CombatComp.Update()

		enemyRect := image.Rect(
			int(enemy.X), int(enemy.Y),
			int(enemy.X)+constants.Tilesize, int(enemy.Y)+constants.Tilesize,
		)

		if enemyRect.Overlaps(playerRect) {
			if enemy.CombatComp.Attack() {
				applyKnockback(enemy.Sprite, player.Sprite, colliders)
				player.CombatComp.Damage(enemy.CombatComp.AttackPower())
			} else {
				pushBack(enemy.Sprite, player.Sprite, colliders)
			}
		}

		for _, anim := range attackAnimations {
			if anim.Finished {
				continue
			}
			animRect := image.Rect(
				int(anim.Position.X)-8, int(anim.Position.Y)-8,
				int(anim.Position.X)+8, int(anim.Position.Y)+8,
			)
			if animRect.Overlaps(enemyRect) {
				enemy.CombatComp.Damage(player.CombatComp.AttackPower())
				anim.Finished = true
			}
		}

		if enemy.CombatComp.Health() <= 0 {
			deadEnemies[idx] = struct{}{}
		}
	}

	newEnemies := make([]*entities.Enemy, 0)
	for idx, e := range enemies {
		if _, dead := deadEnemies[idx]; !dead {
			newEnemies = append(newEnemies, e)
		}
	}
	return newEnemies
}

func applyKnockback(attacker, victim *entities.Sprite, colliders []image.Rectangle) {
	dx := attacker.X - victim.X
	dy := attacker.Y - victim.Y

	length := math.Hypot(dx, dy)
	if length == 0 {
		dx, dy = 1, 0
	} else {
		dx /= length
		dy /= length
	}

	knockbackDistance := 8.0
	attacker.X += dx * knockbackDistance
	attacker.Y += dy * knockbackDistance

	CheckCollisionHorizontal(attacker, colliders)
	CheckCollisionVertical(attacker, colliders)
}

func pushBack(a, b *entities.Sprite, colliders []image.Rectangle) {
	dx := a.X - b.X
	dy := a.Y - b.Y
	length := math.Hypot(dx, dy)
	if length == 0 {
		dx, dy = 1, 0
	} else {
		dx /= length
		dy /= length
	}
	push := 1.0
	a.X += dx * push
	a.Y += dy * push

	CheckCollisionHorizontal(a, colliders)
	CheckCollisionVertical(a, colliders)
}
