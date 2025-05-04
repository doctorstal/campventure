package resources

import (
	"embed"
	"io"
	"log"

	resource "github.com/quasilyte/ebitengine-resource"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

const (
	ImgPlayerSprite resource.ImageID = iota
	ImgGenSky
	ImgGenGround
)

func NewResourceLoader(fs embed.FS, audioContext *audio.Context) *resource.Loader {
	l := resource.NewLoader(audioContext)
	l.OpenAssetFunc = func(path string) io.ReadCloser {
		file, err := fs.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		return file
	}

	l.ImageRegistry.Assign(map[resource.ImageID]resource.ImageInfo{
		ImgPlayerSprite: {Path: "assets/images/Player.png"},
		ImgGenSky:       {Path: "assets/images/generator/simplex_terrain_sky.png"},
		ImgGenGround:    {Path: "assets/images/generator/simplex_terrain_ground.png"},
	})

	return l
}
