package main

import (
	"image/color"

	"github.com/doctorstal/campventure/scenes"
	"github.com/doctorstal/campventure/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type CampVenture struct {
	sceneMap      map[scenes.SceneId]scenes.Scene
	activeSceneId scenes.SceneId
}

// Draw implements ebiten.Game.
func (c *CampVenture) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 66, 99, 255})

	c.sceneMap[c.activeSceneId].Draw(screen)
}

// Layout implements ebiten.Game.
func (c *CampVenture) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 1280, 960
}

// Update implements ebiten.Game.
func (c *CampVenture) Update() error {
	activeScene := c.sceneMap[c.activeSceneId]
	if !activeScene.IsLoaded() {
		activeScene.FirstLoad()
		activeScene.OnEnter()
	}

	nextSceneId := activeScene.Update()
	if nextSceneId == scenes.SceneExit {
		c.sceneMap[c.activeSceneId].OnExit()
		return ebiten.Termination
	}
	if nextSceneId != c.activeSceneId {
		nextScene := c.sceneMap[nextSceneId]
		// if not loaded load scene
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad()
		}
		c.sceneMap[c.activeSceneId].OnExit()
		nextScene.OnEnter()

		c.activeSceneId = nextSceneId
	}

	return nil
}

func NewCampVenture(loader *resource.Loader, mapLoader *tiled.Loader) ebiten.Game {
	cv := &CampVenture{
		sceneMap: map[scenes.SceneId]scenes.Scene{
			scenes.SceneGame:         scenes.NewGameScene(loader, mapLoader),
			scenes.SceneMapGenerator: scenes.NewGeneratorScene(loader),
		},
		activeSceneId: scenes.SceneMapGenerator,
	}
	return cv
}
