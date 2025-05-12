package scenes

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"

	"github.com/doctorstal/campventure/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ojrac/opensimplex-go"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Generator struct {
	w, h                       int
	frequency, fade, threshold float64
	filling                    bool

	groundImage *ebiten.Image
	grassImage  *ebiten.Image

	posX    int
	scrollX int

	noise          opensimplex.Noise
	generatedImage *ebiten.Image
}

func (g *Generator) RandomizeSeed() {
	g.noise = opensimplex.NewNormalized(rand.Int63())
	g.fillImage(g.posX)
}

func (g *Generator) ScrollH(dx int) {

	g.scrollX += dx

	// TODO work on jitter when re-generating terrain
	// Goroutine did not help for some reason, or maybe it's the issue with scroll itself?
	// Maybe do incremental generation instead?

	if !g.filling {
		const scrollBuffer = 100
		if g.scrollX < scrollBuffer {
			g.filling = true
			go func() {
				g.posX -= g.w
				g.fillImage(g.posX)
				g.filling = false
				g.scrollX += g.w
			}()

		} else if g.scrollX > 2*g.w-scrollBuffer {
			g.filling = true
			go func() {
				g.posX += g.w
				g.fillImage(g.posX)
				g.filling = false
				g.scrollX -= g.w
			}()
		}
	}
}

func (g *Generator) Image() *ebiten.Image {
	return g.generatedImage.SubImage(
		image.Rect(g.scrollX, 0, g.scrollX+g.w, g.h),
	).(*ebiten.Image)
}

func (g *Generator) IsSolid(x, y int) bool {
	worldX := g.frequency * float64(x) / float64(g.w)
	worldY := g.frequency * float64(y) / float64(g.w)
	noiseG := g.noise.Eval2(worldX, worldY)

	noiseTransformed := noiseG*(1-g.fade) + g.fade*(float64(y)/float64(g.h))
	return noiseTransformed >= g.threshold
}

func (g *Generator) fillImage(offsetX int) {
	newImage := ebiten.NewImage(3*g.w, g.h)

	gH := g.grassImage.Bounds().Dx()
	gW := g.grassImage.Bounds().Dy()
	giW := g.groundImage.Bounds().Dx()
	drawGrass := func(x int, y int) {
		gX := (x + offsetX) % gW

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y-gH))
		newImage.DrawImage(
			g.grassImage.SubImage(image.Rect(gX, 0, gX+1, gH)).(*ebiten.Image),
			opts,
		)
	}

	for x := range g.w * 3 {
		air := 0

		depth := 100
		for y := range g.h {
			var c color.Color
			if g.IsSolid(x+offsetX, y) {
				depth++
				c = g.groundImage.At((x+offsetX)%giW, y)
			} else {
				if depth > 0 {
					air = 0
				}
				air++
				depth = 0
				continue
			}

			newImage.Set(x, y, c)
			if depth == 20 {
				drawGrass(x, y)
			}
		}
		if depth < 20 && depth > 0 {
			drawGrass(x, g.h+20-depth)
		}
	}

	if g.generatedImage != nil {
		g.generatedImage.Deallocate()
	}

	g.generatedImage = newImage
}

func (g *Generator) FirstLoad() {
	g.scrollX = g.w
	g.fillImage(g.posX)
}

func NewGenerator(loader *resource.Loader, seed int64) *Generator {
	return &Generator{
		w:         500,
		h:         300,
		frequency: 5.5,
		threshold: 0.5,
		fade:      0.5,

		posX: 15000, // far to right, so we can scroll in both directions
		// TODO find how to loop simplex noise
		groundImage: loader.LoadImage(resources.ImgGenGround).Data,
		grassImage:  loader.LoadImage(resources.ImgGenGrass).Data,
		noise:       opensimplex.NewNormalized(seed),
	}
}

type GeneratorScene struct {
	w, h      int
	loader    *resource.Loader
	generator *Generator
	skyImage  *ebiten.Image
	loaded    bool
}

// Draw implements Scene.
func (g *GeneratorScene) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(3.0, 3.0)
	screen.DrawImage(g.skyImage, opts)
	screen.DrawImage(g.generator.Image(), opts)
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %f", ebiten.ActualFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Fade: %f", g.generator.fade))
}

// FirstLoad implements Scene.
func (g *GeneratorScene) FirstLoad() {
	g.loaded = true
	g.generator = NewGenerator(g.loader, 0)
	g.skyImage = g.loader.LoadImage(resources.ImgGenSky).Data
	g.generator.FirstLoad()
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
		g.generator.RandomizeSeed()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.generator.ScrollH(2)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.generator.ScrollH(-2)
	}
	// if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
	// 	g.generator.fade += 0.1
	// 	g.generator.fillImage(g.generator.posX)
	// }
	// if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.generator.fade > -10.0 {
	// 	g.generator.fade -= 0.1
	// 	g.generator.fillImage(g.generator.posX)
	// }

	return SceneMapGenerator
}

func NewGeneratorScene(loader *resource.Loader) Scene {
	return &GeneratorScene{
		loader: loader,
		loaded: false,
	}
}
