package renderer

import (
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	X, Y          float64
	Width, Height int
	Zoom          float64

	shakeX, shakeY float64
}

func (c *Camera) Apply(opts *ebiten.DrawImageOptions) {
	opts.GeoM.Translate(-c.X+c.shakeX, -c.Y+c.shakeY)
	opts.GeoM.Scale(c.Zoom, c.Zoom)
}

func (c *Camera) Shake(intensity float64) {
	c.shakeX = (rand.Float64()*2 - 1) * intensity
	c.shakeY = (rand.Float64()*2 - 1) * intensity
}

func (c *Camera) Update() {
	c.shakeX *= 0.8
	c.shakeY *= 0.8
	if math.Abs(c.shakeX) < 0.1 {
		c.shakeX = 0
	}
	if math.Abs(c.shakeY) < 0.1 {
		c.shakeY = 0
	}
}
