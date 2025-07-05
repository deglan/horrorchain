package systems

import (
	"math"

	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/camera"
	"github.com/deglan/horrorchain/engine/entities"
	"github.com/hajimehoshi/ebiten/v2"
)

func UpdateAttackAnimations(player *entities.Player, attackAnims []*animations.AttackAnimation) []*animations.AttackAnimation {
	newAnims := attackAnims[:0]
	for _, anim := range attackAnims {
		anim.Update()
		if !anim.Finished {
			newAnims = append(newAnims, anim)
		}
	}

	return newAnims
}

func UpdatePlayerAnimation(player *entities.Player) {
	if anim := player.ActiveAnimation(int(player.Dx), int(player.Dy)); anim != nil {
		anim.Update()
	}
}

func CreateAttackAnimation(player *entities.Player, cam *camera.Camera) *animations.AttackAnimation {
	cX, cY := ebiten.CursorPosition()

	playerCenterX := player.X + 8
	playerCenterY := player.Y + 8

	worldClickX := float64(cX) - cam.X
	worldClickY := float64(cY) - cam.Y

	dx := worldClickX - playerCenterX
	dy := worldClickY - playerCenterY
	length := math.Hypot(dx, dy)

	if length == 0 {
		return nil
	}

	direction := struct{ X, Y float64 }{
		X: dx / length,
		Y: dy / length,
	}
	position := struct{ X, Y float64 }{
		X: playerCenterX,
		Y: playerCenterY,
	}

	return animations.NewAttackAnimation(0, 9, 1, 3.0, 1.0, position, direction)
}
