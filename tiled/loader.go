package tiled

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/lafriks/go-tiled"
)

type Loader struct {
	fs fs.FS
}

func (l *Loader) LoadMap(name string) (*TiledMap, error) {
	filePath := path.Join("assets", "maps", name)
	gameMap, err := tiled.LoadFile(filePath, tiled.WithFileSystem(l.fs))
	if err != nil {
		fmt.Printf("error parsing map: %s\n", err.Error())
		return nil, err
	}
	return NewTiledMap(gameMap)
}

func NewLoader(fs fs.FS) *Loader {
	return &Loader{
		fs: fs,
	}
}
