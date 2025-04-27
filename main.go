package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/doctorstal/campventure/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed assets
var assets embed.FS

func main() {
	fmt.Println("Starting camping")
	ebiten.SetWindowSize(1280, 960)

	audioContext := audio.NewContext(44100)
	loader := resources.NewResourceLoader(assets, audioContext)
	game := NewCampVenture(loader)

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
