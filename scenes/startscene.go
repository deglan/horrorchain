package scenes

import (
	"image/color"

	"github.com/deglan/horrorchain/engine/assetloader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type StartScene struct {
	loaded bool
	am     *assetloader.AudioManager
}

func NewStartScene(am *assetloader.AudioManager) *StartScene {
	return &StartScene{
		loaded: false,
		am:     am,
	}
}

func (s *StartScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 0, 0, 255})
	ebitenutil.DebugPrint(screen, "Press enter to start.")
}

func (s *StartScene) FirstLoad() {
	s.loaded = true
}

func (s *StartScene) IsLoaded() bool {
	return s.loaded
}

func (s *StartScene) OnEnter() {
}

func (s *StartScene) OnExit() {
}

func (s *StartScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}

	return StartSceneId
}

var _ Scene = (*StartScene)(nil)
