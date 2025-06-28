package am

import (
	"bytes"
	_ "image/jpeg"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	// Images
	images map[string]*ebiten.Image
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
		images:       make(map[string]*ebiten.Image),
	}
}

func (am *AssetManager) LoadImage(id string, path string) error {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		return err
	}
	am.images[id] = img
	return nil
}

func (am *AssetManager) LoadImageFromBytes(id string, data []byte) error {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	am.images[id] = img
	return nil
}

func (am *AssetManager) GetImage(id string) *ebiten.Image {
	return am.images[id]
}
