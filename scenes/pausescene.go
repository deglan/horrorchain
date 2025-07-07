package scenes

import (
	"image/color"

	"github.com/deglan/horrorchain/engine/assetloader"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PauseScene struct {
	loaded bool
	am     *assetloader.AudioManager
}

func NewPauseScene(am *assetloader.AudioManager) *PauseScene {
	return &PauseScene{
		loaded: false,
		am:     am,
	}
}

func (s *PauseScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 255, 0, 255})
	ebitenutil.DebugPrint(screen, "Press enter to unpause.")
}

func (s *PauseScene) FirstLoad() {
	s.loaded = true
}

func (s *PauseScene) IsLoaded() bool {
	return s.loaded
}

func (s *PauseScene) OnEnter() {
}

func (s *PauseScene) OnExit() {
}

func (s *PauseScene) Update() SceneId {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return GameSceneId
	}

	return PauseSceneId
}

var _ Scene = (*PauseScene)(nil)
