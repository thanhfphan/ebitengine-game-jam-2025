package renderer

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/thanhfphan/ebitengj2025/internal/am"
	"github.com/thanhfphan/ebitengj2025/internal/world"
)

type Renderer struct {
	Camera *Camera
	AM     *am.AssetManager

	Background color.Color
}

func New(am *am.AssetManager) *Renderer {
	cam := &Camera{
		X:      0,
		Y:      0,
		Width:  1280,
		Height: 720,
	}
	return &Renderer{
		Camera:     cam,
		Background: color.RGBA{0, 0, 0, 255},
	}
}

func (r *Renderer) Update() {
}

func (r *Renderer) DrawWorld(screen *ebiten.Image, world *world.World) {
	opts := &ebiten.DrawImageOptions{}
	r.Camera.Apply(opts)

	// for _, c := range world.Cards {
	// 	cardOpts := &ebiten.DrawImageOptions{}
	// 	r.Camera.Apply(cardOpts)
	// 	cardOpts.GeoM.Translate(c.X, c.Y)
	// screen.DrawImage(r.AM.Image(c.SpriteID), cardOpts)
	// }
}

func (r *Renderer) SetCameraSize(width, height int) {
	r.Camera.Width = width
	r.Camera.Height = height
}

func (r *Renderer) BackgroundColor() color.Color {
	return r.Background
}

func (r *Renderer) DrawDebug(screen *ebiten.Image, text string) {
	ebitenutil.DebugPrint(screen, text)
}

func (r *Renderer) DrawDebugAt(screen *ebiten.Image, text string, x, y int) {
	ebitenutil.DebugPrintAt(screen, text, x, y)
}
