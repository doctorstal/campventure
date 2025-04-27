package main

import (
	"fmt"
	"image/color"

	"github.com/doctorstal/campventure/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
)

type CampVenture struct {
	Player *entities.Player
}

// Draw implements ebiten.Game.
func (c *CampVenture) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{33, 66, 99, 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Vertical speed: %f", c.Player.Dy))

	opts := c.Player.DrawOptions()
	opts.GeoM.Translate(c.Player.Bounds().Lo().X, c.Player.Bounds().Lo().Y)
	opts.GeoM.Scale(4, 4)
	screen.DrawImage(c.Player.Img(), opts)
}

// Layout implements ebiten.Game.
func (c *CampVenture) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 1280, 960
}

// Update implements ebiten.Game.
func (c *CampVenture) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		c.Player.GoRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		c.Player.GoLeft()
	}
	playerFalls := c.Player.Bounds().Hi().Y < 240
	if playerFalls {
		c.Player.Dy += 0.3
	} else if c.Player.Dy > 0 {
		c.Player.Dy = 0
	}

	if !playerFalls && ebiten.IsKeyPressed(ebiten.KeySpace) {
		c.Player.Jump()
	}

	c.Player.Update()
	return nil
}

func NewCampVenture(loader *resource.Loader) ebiten.Game {
	return &CampVenture{
		Player: entities.NewPlayer(loader).(*entities.Player),
	}
}
