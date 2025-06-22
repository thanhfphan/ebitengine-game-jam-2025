package am

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type fontKey struct {
	Name string
	Size float64
}

func (am *AssetManager) LoadFont(name string, data []byte, size float64) error {
	if am.fonts == nil {
		am.fonts = make(map[fontKey]font.Face)
	}
	if am.fontData == nil {
		am.fontData = make(map[string][]byte)
	}

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

	key := fontKey{name, size}
	am.fonts[key] = face
	am.fontData[name] = data
	return nil
}

func (am *AssetManager) GetFont(name string, size float64) font.Face {
	key := fontKey{name, size}
	if face, ok := am.fonts[key]; ok {
		return face
	}

	data, ok := am.fontData[name]
	if !ok {
		panic(fmt.Sprintf("font data for '%s' not loaded", name))
	}

	if err := am.LoadFont(name, data, size); err != nil {
		panic(fmt.Sprintf("failed to load font '%s' size %.1f: %v", name, size, err))
	}
	return am.fonts[key]
}
