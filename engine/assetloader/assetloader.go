package assetloader

import (
	"log"

	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/components"
	"github.com/deglan/horrorchain/engine/entities"
	"github.com/deglan/horrorchain/engine/systems"
	"github.com/deglan/horrorchain/engine/tile"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("failed to load image %s: %v", path, err)
	}
	return img
}

func LoadTilemap(path string) (*tile.TilemapJSON, []tile.Tileset) {
	tmap, err := tile.NewTilemapJSON(path)
	if err != nil {
		log.Fatalf("failed to load tilemap %s: %v", path, err)
	}
	tilesets, err := tmap.GenTilesets()
	if err != nil {
		log.Fatalf("failed to generate tilesets: %v", err)
	}
	return tmap, tilesets
}

func LoadPlayer(img *ebiten.Image) *entities.Player {
	return &entities.Player{
		Sprite: &entities.Sprite{
			Img: img,
			X:   50.0,
			Y:   50.0,
		},
		Health: 9,
		Animations: map[entities.PlayerState]*animations.Animation{
			entities.Up:    animations.NewAnimation(5, 13, 4, 20.0, 1.0),
			entities.Down:  animations.NewAnimation(4, 12, 4, 20.0, 1.0),
			entities.Left:  animations.NewAnimation(6, 14, 4, 20.0, 1.0),
			entities.Right: animations.NewAnimation(7, 15, 4, 20.0, 1.0),
		},
		CombatComp: components.NewBasicCombat(9, 1),
	}
}

func LoadEnemies(img *ebiten.Image) []*entities.Enemy {
	return []*entities.Enemy{
		{
			Sprite: &entities.Sprite{
				Img: img,
				X:   100.0,
				Y:   100.0,
			},
			FollowsPlayer: true,
			CombatComp:    components.NewEnemyCombat(3, 1, 30),
		},
		{
			Sprite: &entities.Sprite{
				Img: img,
				X:   150.0,
				Y:   50.0,
			},
			FollowsPlayer: true,
			CombatComp:    components.NewEnemyCombat(3, 1, 30),
		},
	}
}

func LoadPotions(img *ebiten.Image) []*entities.Potion {
	return []*entities.Potion{
		{
			Sprite: &entities.Sprite{
				Img: img,
				X:   210.0,
				Y:   100.0,
			},
			AmtHeal: 1.0,
		},
	}
}

func DrawHearts(screen *ebiten.Image, img *ebiten.Image, anim *animations.Animation, hp int) {
	for i := 0; i < hp; i++ {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(16*i), 16) // lewy górny róg, z odstępem

		frame := anim.Frame()
		screen.DrawImage(
			img.SubImage(systems.GetHeartFrame(img, frame)).(*ebiten.Image),
			opts,
		)
	}
}
