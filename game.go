package main

import (
	"github.com/deglan/horrorchain/engine/assetloader"
	"github.com/deglan/horrorchain/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneMap      map[scenes.SceneId]scenes.Scene
	activeSceneId scenes.SceneId
	audioManager  *assetloader.AudioManager
}

func NewGame() *Game {
	audioManager := assetloader.NewAudioManager()
	err := audioManager.LoadSound("attack", "assets/sounds/effects/attack_effect.wav")
	if err != nil {
		panic(err)
	}
	sceneMap := map[scenes.SceneId]scenes.Scene{
		scenes.GameSceneId:     scenes.NewGameScene(audioManager),
		scenes.StartSceneId:    scenes.NewStartScene(audioManager),
		scenes.PauseSceneId:    scenes.NewPauseScene(audioManager),
		scenes.GameOverSceneId: scenes.NewGameOverScene(audioManager),
	}
	activeSceneId := scenes.StartSceneId
	sceneMap[activeSceneId].FirstLoad()
	return &Game{
		sceneMap,
		activeSceneId,
		audioManager,
	}
}

func (g *Game) Update() error {
	nextSceneId := g.sceneMap[g.activeSceneId].Update()
	// switched scenes
	if nextSceneId == scenes.ExitSceneId {
		g.sceneMap[g.activeSceneId].OnExit()
		return ebiten.Termination
	}
	if nextSceneId != g.activeSceneId {
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded? then load in
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad()
		}
		nextScene.OnEnter()
		g.sceneMap[g.activeSceneId].OnExit()
	}
	g.activeSceneId = nextSceneId
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
