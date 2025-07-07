package scenes

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/deglan/horrorchain/constants"
	"github.com/deglan/horrorchain/engine/animations"
	"github.com/deglan/horrorchain/engine/assetloader"
	"github.com/deglan/horrorchain/engine/camera"
	"github.com/deglan/horrorchain/engine/entities"
	"github.com/deglan/horrorchain/engine/spritesheet"
	"github.com/deglan/horrorchain/engine/systems"
	"github.com/deglan/horrorchain/engine/tile"
)

type GameScene struct {
	loaded            bool
	player            *entities.Player
	playerSpriteSheet *spritesheet.SpriteSheet
	enemies           []*entities.Enemy
	potions           []*entities.Potion
	tilemapJSON       *tile.TilemapJSON
	tilesets          []tile.Tileset
	tilemapImg        *ebiten.Image
	camera            *camera.Camera
	colliders         []image.Rectangle
	attackAnimations  []*animations.AttackAnimation
	healthImg         *ebiten.Image
	healthAnim        *animations.Animation
	am                *assetloader.AudioManager
}

func NewGameScene(am *assetloader.AudioManager) *GameScene {
	return &GameScene{
		player:            nil,
		playerSpriteSheet: nil,
		enemies:           make([]*entities.Enemy, 0),
		potions:           make([]*entities.Potion, 0),
		tilemapJSON:       nil,
		tilesets:          nil,
		tilemapImg:        nil,
		camera:            nil,
		colliders:         make([]image.Rectangle, 0),
		loaded:            false,
		attackAnimations:  make([]*animations.AttackAnimation, 0),
		healthImg:         nil,
		healthAnim:        nil,
		am:                am,
	}
}

func (g *GameScene) IsLoaded() bool {
	return g.loaded
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	systems.DrawTileLayers(screen, g.tilemapJSON, g.tilesets, g.camera.X, g.camera.Y)
	systems.DrawPlayer(screen, g.player, g.playerSpriteSheet, g.camera.X, g.camera.Y)
	systems.DrawEnemies(screen, g.enemies, g.camera.X, g.camera.Y)
	systems.DrawPotions(screen, g.potions, g.camera.X, g.camera.Y)
	systems.DrawColliders(screen, g.colliders, g.camera.X, g.camera.Y)
	systems.DrawAttackAnimations(screen, g.attackAnimations, g.camera.X, g.camera.Y)

	systems.DrawHearts(screen, g.healthImg, g.healthAnim, g.player.CombatComp.Health())
	systems.DrawEnemyDebugRects(screen, g.enemies, g.camera.X, g.camera.Y)

}

func (g *GameScene) FirstLoad() {
	playerImg := assetloader.LoadImage("assets/images/ninja.png")
	skeletonImg := assetloader.LoadImage("assets/images/skeleton.png")
	potionImg := assetloader.LoadImage("assets/images/potion.png")
	tilemapImg := assetloader.LoadImage("assets/images/atlas_floor.png")
	g.healthImg = assetloader.LoadImage("assets/images/HealthUI.png")
	g.healthAnim = animations.NewAnimation(0, 5, 1, 4.0, 4.0)

	tilemapJSON, tilesets := assetloader.LoadTilemap("assets/maps/firstRoom.json")

	g.player = assetloader.LoadPlayer(playerImg)
	g.playerSpriteSheet = spritesheet.NewSpriteSheet(4, 7, constants.Tilesize)
	g.enemies = assetloader.LoadEnemies(skeletonImg)
	g.potions = assetloader.LoadPotions(potionImg)
	g.tilemapJSON = tilemapJSON
	g.tilemapImg = tilemapImg
	g.tilesets = tilesets
	g.camera = camera.NewCamera(0.0, 0.0)

	g.colliders = append(g.colliders, g.tilemapJSON.ExtractCollidersFromLayer("colliders")...)
	g.colliders = append(g.colliders, g.tilemapJSON.ExtractCollidersFromLayer("walls")...)

	g.healthAnim = animations.NewAnimation(0, 5, 1, 4.0, 4.0)
	g.loaded = true
}

func (g *GameScene) OnEnter() {
}

func (g *GameScene) OnExit() {
}

// Update implements Scene.
func (g *GameScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ExitSceneId
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		return PauseSceneId
	}
	if g.player.CombatComp.Health() <= 0 {
		return GameOverSceneId
	}

	g.attackAnimations = systems.UpdateAttackAnimations(g.player, g.attackAnimations)
	wallsGrid := g.tilemapJSON.BuildGrid("walls", "colliders")
	systems.MovePlayer(g.player, g.colliders)
	systems.MoveEnemies(g.enemies, g.player, g.colliders, wallsGrid)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		if anim := systems.CreateAttackAnimation(g.player, g.camera); anim != nil {
			g.am.PlayWithVolume("attack", 0.1)
			g.attackAnimations = append(g.attackAnimations, anim)
		}
	}

	g.enemies = systems.HandleCombat(g.player, g.enemies, g.attackAnimations, g.colliders)
	systems.UpdatePlayerAnimation(g.player)
	g.potions = systems.CollectPotions(g.player, g.potions)

	g.camera.FollowTarget(g.player.X, g.player.Y, constants.WindowSizeWidth, constants.WindowSizeHeight)
	g.camera.Constrain(
		float64(g.tilemapJSON.Layers[0].Width)*constants.Tilesize,
		float64(g.tilemapJSON.Layers[0].Height)*constants.Tilesize,
		constants.WindowSizeWidth,
		constants.WindowSizeHeight,
	)

	return GameSceneId
}

var _ Scene = (*GameScene)(nil)
