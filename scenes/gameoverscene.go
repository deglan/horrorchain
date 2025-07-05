package scenes

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type GameOverScene struct {
	loaded bool
}

func NewGameOverScene() *GameOverScene {
	return &GameOverScene{
		loaded: false,
	}
}

func (s *GameOverScene) Draw(screen *ebiten.Image) {

	screen.Fill(color.Black)
	msg := "Game Over!"

	width, height := text.Measure(msg, DefaultFontFace, 4)

	screenW, screenH := 320, 240

	x := (float64(screenW) - width) / 2
	y := (float64(screenH) - height) / 2

	opts := &text.DrawOptions{}
	opts.GeoM.Translate(x, y)
	opts.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, msg, DefaultFontFace, opts)
}

func (s *GameOverScene) FirstLoad() {
	if DefaultFontFace == nil {
		LoadFont()
	}
	s.loaded = true
}

func (s *GameOverScene) IsLoaded() bool {
	return s.loaded
}

func (s *GameOverScene) OnEnter() {
}

func (s *GameOverScene) OnExit() {
}

func (s *GameOverScene) Update() SceneId {
	return GameOverSceneId
}

var _ Scene = (*GameOverScene)(nil)
