package scenes

import (
	"fmt"

	"github.com/doctorstal/campventure/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
)

type GameScene struct {
	player   *entities.Player
	isLoaded bool
}

// Draw implements Scene.
func (g *GameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Vertical speed: %f", g.player.Dy))

	opts := g.player.DrawOptions()
	opts.GeoM.Translate(g.player.Bounds().Lo().X, g.player.Bounds().Lo().Y)
	opts.GeoM.Scale(4, 4)
	screen.DrawImage(g.player.Img(), opts)
}

// FirstLoad implements Scene.
func (g *GameScene) FirstLoad() {
	g.isLoaded = true
}

// IsLoaded implements Scene.
func (g *GameScene) IsLoaded() bool {
	return g.isLoaded
}

// OnEnter implements Scene.
func (g *GameScene) OnEnter() {
}

// OnExit implements Scene.
func (g *GameScene) OnExit() {
}

// Update implements Scene.
func (g *GameScene) Update() SceneId {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return SceneExit
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.GoRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.GoLeft()
	}
	playerFalls := g.player.Bounds().Hi().Y < 240
	if playerFalls {
		g.player.Dy += 0.3
	} else if g.player.Dy > 0 {
		g.player.Dy = 0
	}

	if !playerFalls && ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.player.Jump()
	}

	g.player.Update()

	return SceneGame
}

func NewGameScene(loader *resource.Loader) Scene {
	return &GameScene{
		player: entities.NewPlayer(loader).(*entities.Player),
	}
}
