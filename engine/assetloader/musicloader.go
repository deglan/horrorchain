package assetloader

import (
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type AudioManager struct {
	ctx       *audio.Context
	rawSounds map[string]*os.File
	volume    float64
}

func NewAudioManager() *AudioManager {
	ctx := audio.NewContext(44100)
	return &AudioManager{
		ctx:       ctx,
		rawSounds: make(map[string]*os.File),
		volume:    1.0,
	}
}

func (am *AudioManager) LoadSound(name, path string) error {
	sound, err := os.Open(path)
	if err != nil {
		return err
	}
	am.rawSounds[name] = sound
	return nil
}

func (am *AudioManager) Play(name string) {
	am.PlayWithVolume(name, am.volume)
}

func (am *AudioManager) PlayWithVolume(name string, volume float64) {
	raw, ok := am.rawSounds[name]
	if !ok {
		return
	}
	raw.Seek(0, io.SeekStart)

	stream, err := wav.DecodeF32(raw)
	if err != nil {
		return
	}

	player, err := am.ctx.NewPlayer(stream)
	if err != nil {
		return
	}

	player.SetVolume(volume)
	player.Play()
}
