package scenes

import (
	"image"
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
	grassImage  *ebiten.Image

	noise          opensimplex.Noise
	generatedImage *ebiten.Image
	loaded         bool
}

// Draw implements Scene.
func (g *GeneratorScene) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(2.0, 2.0)
	screen.DrawImage(g.generatedImage, opts)
}

// FirstLoad implements Scene.
func (g *GeneratorScene) FirstLoad() {
	g.loaded = true

	g.generatedImage = ebiten.NewImage(g.w, g.h)

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
				c = g.groundImage.At(x, y)

			} else {
				depth = 0
				c = g.skyImage.At(x, y)
			}

			g.generatedImage.Set(x, y, c)
			if depth == 10 {
				gH := g.grassImage.Bounds().Dx()
				gW := g.grassImage.Bounds().Dy()
				gX := x % gW

				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x), float64(y-gH))
				g.generatedImage.DrawImage(
					g.grassImage.SubImage(image.Rect(gX, 0, gX+1, gH)).(*ebiten.Image),
					opts,
				)
			}
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
		grassImage:  loader.LoadImage(resources.ImgGenGrass).Data,
		loader:      loader,
		noise:       opensimplex.NewNormalized(2),
	}
}
