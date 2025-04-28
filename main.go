package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/doctorstal/campventure/resources"
	"github.com/doctorstal/campventure/tiled"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//go:embed assets
var assets embed.FS

//go:embed assets/maps
var maps embed.FS

func main() {
	fmt.Println("Starting camping")
	ebiten.SetWindowSize(1280, 960)

	audioContext := audio.NewContext(44100)
	loader := resources.NewResourceLoader(assets, audioContext)
	mapLoader := tiled.NewLoader(maps)
	game := NewCampVenture(loader, mapLoader)

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
