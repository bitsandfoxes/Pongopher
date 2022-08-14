package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type GameObject struct {
	image    *ebiten.Image
	size     Vector
	position Vector
	speed    float64
	bounds   Bounds
}

type Bounds struct {
	left   float64
	right  float64
	top    float64
	bottom float64
}

func MakeGameObject(width, height int, x, y, speed float64, color color.RGBA) GameObject {
	gameObject := GameObject{
		image: ebiten.NewImage(width, height),
		size:  Vector{float64(width), float64(height)},
		speed: speed,
	}

	gameObject.SetPosition(Vector{x, y})
	gameObject.image.Fill(color)

	return gameObject
}

func (o *GameObject) SetPosition(v Vector) {
	o.position = v
	o.bounds = Bounds{
		left:   v.x - float64(o.size.x)/2,
		right:  v.x + float64(o.size.x)/2,
		top:    v.y - float64(o.size.y)/2,
		bottom: v.y + float64(o.size.y)/2,
	}
}

func (o *GameObject) GetScreenPosition() (float64, float64) {
	return o.bounds.left, o.bounds.top
}

func (o *GameObject) Hit(other GameObject) bool {
	if o.bounds.right >= other.bounds.left &&
		o.bounds.left <= other.bounds.right &&
		o.bounds.top <= other.bounds.bottom &&
		o.bounds.bottom >= other.bounds.top {
		return true
	}

	return false
}
