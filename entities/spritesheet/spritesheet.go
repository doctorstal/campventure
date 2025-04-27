package spritesheet

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteSheet struct {
	img           *ebiten.Image
	WidthInTiles  int
	HeightInTiles int
	TileSize      int
}

func (s *SpriteSheet) Rect(index int) image.Rectangle {
	x := index % s.WidthInTiles * s.TileSize
	y := index / s.WidthInTiles * s.TileSize
	return image.Rect(x, y, x+s.TileSize, y+s.TileSize)
}

func (s *SpriteSheet) Img(index int) *ebiten.Image {
	return s.img.SubImage(s.Rect(index)).(*ebiten.Image)
}

// Creates new sprite sheet, takes width and height in tiles and tile size
func NewSpriteSheet(img *ebiten.Image, w, h, t int) *SpriteSheet {
	return &SpriteSheet{
		img,
		w,
		h,
		t,
	}
}
