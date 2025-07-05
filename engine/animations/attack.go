package animations

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type AttackAnimation struct {
	*Animation
	Frames    []*ebiten.Image
	Finished  bool
	Position  struct{ X, Y float64 }
	Direction struct{ X, Y float64 }
}

func NewAttackAnimation(first, last, step int, speed, frameCounter float32, position, direction struct{ X, Y float64 }) *AttackAnimation {
	return &AttackAnimation{
		NewAnimation(first, last, step, speed, frameCounter),
		createFrameList(),
		false,
		position,
		direction,
	}
}

func (a *AttackAnimation) Update() {
	if a.Finished {
		return
	}

	a.frameCounter--
	if a.frameCounter <= 0 {
		a.frameCounter = a.SpeedInTps
		a.frame += a.Step
		if a.frame > a.Last || a.frame >= len(a.Frames) {
			a.Finished = true
			return
		}
	}

	const speed = 2.0
	a.Position.X += a.Direction.X * speed
	a.Position.Y += a.Direction.Y * speed
}

func (a *AttackAnimation) Draw(screen *ebiten.Image, cameraX, cameraY float64) {
	if a.Finished || a.frame >= len(a.Frames) {
		return
	}

	img := a.Frames[a.frame]
	opts := &ebiten.DrawImageOptions{}

	scale := 0.125
	opts.GeoM.Scale(scale, scale)

	opts.GeoM.Translate(-128*scale, -128*scale)

	opts.GeoM.Translate(a.Position.X+cameraX, a.Position.Y+cameraY)

	screen.DrawImage(img, opts)
}

func createFrameList() []*ebiten.Image {
	frames := make([]*ebiten.Image, 0, 10)
	for i := 0; i < 10; i++ {
		attackFrame, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("assets/images/attack/Explosion_blue_circle/Explosion_blue_circle%d.png", i+1))
		if err != nil {
			panic(err)
		}
		frames = append(frames, attackFrame)
	}
	return frames
}
