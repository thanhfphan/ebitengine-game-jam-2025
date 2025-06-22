package am

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type AssetManager struct {
	fonts    map[fontKey]font.Face
	fontData map[string][]byte // map font name → raw []byte (để reuse khi load size mới)
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		fonts:    make(map[fontKey]font.Face),
		fontData: make(map[string][]byte),
	}
}

func (am *AssetManager) GetCardImage(cardID string) *ebiten.Image {
	// TODO: Implement actual image loading
	img := ebiten.NewImage(80, 120)
	return img
}

func (am *AssetManager) GetCardBackImage() *ebiten.Image {
	// TODO: Implement actual card back image loading
	img := ebiten.NewImage(80, 120)
	return img
}
