package scenes

import "github.com/hajimehoshi/ebiten/v2"

type SceneId int

const (
	SceneIntro SceneId = iota
	SceneGame
	ScenePause
	SceneGameOver
	SceneMapGenerator
	SceneExit
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update() SceneId
	FirstLoad()
	OnEnter()
	OnExit()
	IsLoaded() bool
}
