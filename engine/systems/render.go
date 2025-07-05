package systems

import (
	"fmt"
	"image"
	"image/color"

	"github.com/deglan/horrorchain/constants"
	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/entities"
	"github.com/deglan/horrorchain/engine/spritesheet"
	"github.com/deglan/horrorchain/engine/tile"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawTileLayers(screen *ebiten.Image, tilemapJSON *tile.TilemapJSON, tilesets []tile.Tileset, cameraX, cameraY float64) {
	for _, layer := range tilemapJSON.Layers {
		for index, id := range layer.Data {
			if id == 0 {
				continue
			}
			x := index % layer.Width * constants.Tilesize
			y := index / layer.Width * constants.Tilesize

			var img *ebiten.Image
			for _, ts := range tilesets {
				if ts.Contains(id) {
					img = ts.Img(id)
					break
				}
			}
			if img == nil {
				continue
			}

			opts := ebiten.DrawImageOptions{}
			imgHeight := img.Bounds().Dy()
			tileOffsetY := constants.Tilesize - imgHeight
			opts.GeoM.Translate(float64(x), float64(y+tileOffsetY))

			opts.GeoM.Translate(cameraX, cameraY)

			fmt.Printf("Layer: %s, TileID: %d, Using tileset: %T\n", layer.Name, id, tilemapJSON.Tilesets)

			screen.DrawImage(img, &opts)
		}
	}
}

func DrawPlayer(screen *ebiten.Image, player *entities.Player, playerSpriteSheet *spritesheet.SpriteSheet, cameraX, cameraY float64) {
	if player.Img == nil {
		return
	}

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(player.X+cameraX, player.Y+cameraY)

	frame := 0
	if anim := player.ActiveAnimation(int(player.Dx), int(player.Dy)); anim != nil {
		frame = anim.Frame()
	}

	screen.DrawImage(
		player.Img.SubImage(playerSpriteSheet.Rect(frame)).(*ebiten.Image),
		&opts,
	)
}

func DrawHearts(screen *ebiten.Image, img *ebiten.Image, anim *animations.Animation, hp int) {
	for i := 0; i < hp; i++ {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(16*i), 16) // lewy górny róg, z odstępem

		frame := anim.Frame()
		screen.DrawImage(
			img.SubImage(GetHeartFrame(img, frame)).(*ebiten.Image),
			opts,
		)
	}
}

func GetHeartFrame(img *ebiten.Image, frame int) image.Rectangle {
	const width, height = 11, 11
	y := frame * height
	return image.Rect(0, y, width, y+height)
}

func DrawEnemies(screen *ebiten.Image, enemies []*entities.Enemy, cameraX, cameraY float64) {
	opts := ebiten.DrawImageOptions{}

	for _, enemy := range enemies {
		opts.GeoM.Translate(enemy.X, enemy.Y)
		opts.GeoM.Translate(cameraX, cameraY)

		screen.DrawImage(
			enemy.Img.SubImage(image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
}

func DrawPotions(screen *ebiten.Image, potions []*entities.Potion, cameraX, cameraY float64) {
	opts := ebiten.DrawImageOptions{}

	for _, potion := range potions {
		opts.GeoM.Translate(potion.X, potion.Y)
		opts.GeoM.Translate(cameraX, cameraY)

		screen.DrawImage(
			potion.Img.SubImage(image.Rect(0, 0, constants.Tilesize, constants.Tilesize)).(*ebiten.Image),
			&opts,
		)
		opts.GeoM.Reset()
	}
}

func DrawColliders(screen *ebiten.Image, colliders []image.Rectangle, cameraX, cameraY float64) {
	for _, collider := range colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(cameraX),
			float32(collider.Min.Y)+float32(cameraY),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0, color.NRGBA{0, 0, 0, 0}, true,
		)
	}
}

func DrawAttackAnimations(screen *ebiten.Image, animations []*animations.AttackAnimation, cameraX, cameraY float64) {
	for _, anim := range animations {
		anim.Draw(screen, cameraX, cameraY)
	}
}
