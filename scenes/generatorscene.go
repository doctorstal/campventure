package scenes

import (
	"image/color"
	"math/rand"

	"github.com/doctorstal/campventure/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ojrac/opensimplex-go"
	resource "github.com/quasilyte/ebitengine-resource"
)

type GeneratorScene struct {
	w, h   int
	loader *resource.Loader

	groundImage *ebiten.Image
	skyImage    *ebiten.Image

	noise      opensimplex.Noise
	noiseImage *ebiten.Image
	loaded     bool
}

// Draw implements Scene.
func (g *GeneratorScene) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.noiseImage, nil)
}

// FirstLoad implements Scene.
func (g *GeneratorScene) FirstLoad() {
	g.loaded = true

	g.noiseImage = ebiten.NewImage(g.w, g.h)

	g.fillImage()
}

func (g *GeneratorScene) fillImage() {
	frequency := 3.5
	threshold := 0.5
	fade := 0.5

	isSolid := func(x, y int) bool {
		noiseG := g.noise.Eval2(frequency*float64(x)/float64(g.w), frequency*float64(y)/float64(g.w))
		return noiseG*(1-fade)+fade*(float64(y)/float64(g.h)) >= threshold
	}

	for x := range g.w {
		depth := 100
		for y := range g.h {
			var c color.Color
			if isSolid(x, y) {
				depth++
				if depth <= 10 {
					c = color.RGBA{33, 99, 66, 1}
				} else {
					c = g.groundImage.At(x, y)
				}
			} else {
				depth = 0
				c = g.skyImage.At(x, y)
			}

			g.noiseImage.Set(x, y, c)
		}
	}
}

// IsLoaded implements Scene.
func (g *GeneratorScene) IsLoaded() bool {
	return g.loaded
}

// OnEnter implements Scene.
func (g *GeneratorScene) OnEnter() {
}

// OnExit implements Scene.
func (g *GeneratorScene) OnExit() {
}

// Update implements Scene.
func (g *GeneratorScene) Update() SceneId {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		g.noise = opensimplex.NewNormalized(rand.Int63())

		g.fillImage()
	}

	return SceneMapGenerator
}

func NewGeneratorScene(loader *resource.Loader) Scene {
	return &GeneratorScene{
		w:           500,
		h:           300,
		skyImage:    loader.LoadImage(resources.ImgGenSky).Data,
		groundImage: loader.LoadImage(resources.ImgGenGround).Data,
		loader:      loader,
		noise:       opensimplex.NewNormalized(2),
	}
}
