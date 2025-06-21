package am

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type AssetManager struct {
	fonts map[string]font.Face
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		fonts: make(map[string]font.Face),
	}
}

func (am *AssetManager) LoadFont(name string, data []byte, size float64) error {
	ttf, err := opentype.Parse(data)
	if err != nil {
		return fmt.Errorf("error parsing font: %v", err)
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("error creating font face: %v", err)
	}

	am.fonts[name] = face

	return nil
}

func (am *AssetManager) GetFont(name string) font.Face {
	return am.fonts[name]
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
