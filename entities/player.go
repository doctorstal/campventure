package entities

import (
	"math"

	"github.com/golang/geo/r2"

	"github.com/doctorstal/campventure/animations"
	"github.com/doctorstal/campventure/entities/spritesheet"
	"github.com/doctorstal/campventure/resources"
	"github.com/hajimehoshi/ebiten/v2"
	resource "github.com/quasilyte/ebitengine-resource"
)

type Animator interface {
	// Returns true if it should be removed
	Update() bool
	Bounds() r2.Rect
	Z() int
	Img() *ebiten.Image
	DrawOptions() *ebiten.DrawImageOptions
}

type SpriteState int

const (
	StateIdle SpriteState = iota
	StateRun
	StateJump
)

type PlayerDirection int

const (
	DirecitonRight PlayerDirection = iota
	DirecitonLeft
)

type Sprite struct {
	spriteSheet         *spritesheet.SpriteSheet
	animation           animations.Animation
	state               SpriteState
	x, y, width, height float64
}

type Player struct {
	*Sprite
	Dx, Dy     float64
	direction  PlayerDirection
	animations map[SpriteState]animations.Animation
}

func (p *Player) GoRight() {
	if p.state != StateJump {
		p.direction = DirecitonRight
		p.state = StateRun
		p.Dx = 2
	}
}

func (p *Player) GoLeft() {
	if p.state != StateJump {
		p.direction = DirecitonLeft
		p.state = StateRun
		p.Dx = -2
	}
}

func (p *Player) Jump() {
	if p.state == StateJump {
		return
	}
	lastState := p.state
	lastDx := p.Dx
	p.Dx = 0

	p.state = StateJump
	// p.animation = animations.NewOneTimeAnimation(0, 7, 1, 10.0, false)
	p.animations[StateJump] = animations.NewCallBackAnimation(0, 7, 1, 10.0, func(frame int) bool {
		if frame == 2 {
			p.Dy = -5
			p.Dx = lastDx
		}
		if frame == 5 {
			p.Dx = 0
		}
		if frame == 6 {
			p.state = lastState
			return true
		}
		return false
	})
}

// DrawOptions implements Animator.
func (p *Player) DrawOptions() *ebiten.DrawImageOptions {
	// Flip
	opts := ebiten.DrawImageOptions{}
	if p.direction == DirecitonLeft {
		opts.GeoM.Translate(-p.width/2, 0)
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(p.width/2, 0)
	}
	return &opts
}

// Img implements Animator.
func (p *Player) Img() *ebiten.Image {
	return p.spriteSheet.Img(p.animation.Frame())
}

// Z implements Animator.
func (p *Player) Z() int {
	return 0
}

// Bounds implements Animator.
func (p *Player) Bounds() r2.Rect {
	return r2.RectFromCenterSize(
		r2.Point{X: p.x, Y: p.y},
		r2.Point{X: p.width, Y: p.height},
	)
}

func (p *Player) Update() bool {
	p.Move()
	p.animation = p.animations[p.state]
	if p.state == StateRun && math.Abs(p.Dx) < 0.1 {
		p.state = StateIdle
	}
	p.animation.Update()
	return false // Never remove player
}

func (p *Player) Move() {
	p.x += p.Dx
	p.y += p.Dy
	p.Dy *= 0.9
	if math.Abs(p.Dy) < 0.1 {
		p.Dy = 0
		p.Dx *= 0.7
	}
}

func NewPlayer(loader *resource.Loader) Animator {
	playerImg := loader.LoadImage(resources.ImgPlayerSprite).Data
	spriteSheet := spritesheet.NewSpriteSheet(playerImg, 6, 2, 32)
	return &Player{
		Sprite: &Sprite{
			spriteSheet: spriteSheet,
			animation:   animations.NewLoopAnimation(7, 10, 1, 10.0),
			x:           0,
			y:           0,
			width:       32,
			height:      32,
		},
		animations: map[SpriteState]animations.Animation{
			StateIdle: animations.NewSingleFrameAnimation(0),
			StateRun:  animations.NewLoopAnimation(7, 10, 1, 10.0),
		},
	}
}
