package am

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"golang.org/x/image/font"
)

type AssetManager struct {
	fonts    map[fontKey]font.Face
	fontData map[string][]byte // map font name → raw []byte (để reuse khi load size mới)

	// Audio context and players
	audioContext *audio.Context
	sounds       map[string]*Sound
	music        map[string]*Music

	// Volume settings
	masterVolume float64
	musicVolume  float64
}

func NewAssetManager() *AssetManager {
	// Create audio context with 44.1 kHz sample rate
	audioContext := audio.NewContext(44100)

	return &AssetManager{
		fonts:        make(map[fontKey]font.Face),
		fontData:     make(map[string][]byte),
		audioContext: audioContext,
		sounds:       make(map[string]*Sound),
		music:        make(map[string]*Music),
		masterVolume: 1.0, // Default to full volume
		musicVolume:  1.0, // Default to full volume
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
