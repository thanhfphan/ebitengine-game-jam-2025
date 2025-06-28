package am

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Sound represents a short sound effect
type Sound struct {
	audioContext *audio.Context
	data         []byte
	format       string
	volume       float64
	am           *AssetManager
}

// Play plays the sound effect once
func (s *Sound) Play() error {
	var player *audio.Player
	var err error

	switch s.format {
	case "wav":
		wavReader, errwav := wav.DecodeWithSampleRate(s.audioContext.SampleRate(), bytes.NewReader(s.data))
		if errwav != nil {
			return fmt.Errorf("error decoding WAV: %v", errwav)
		}
		player, err = s.audioContext.NewPlayer(wavReader)
	case "ogg":
		oggReader, errr := vorbis.DecodeWithSampleRate(s.audioContext.SampleRate(), bytes.NewReader(s.data))
		if errr != nil {
			return fmt.Errorf("error decoding OGG: %v", errr)
		}
		player, err = s.audioContext.NewPlayer(oggReader)
	default:
		return fmt.Errorf("unsupported audio format: %s", s.format)
	}

	if err != nil {
		return fmt.Errorf("error creating audio player: %v", err)
	}

	// Set volume based on master volume
	player.SetVolume(s.volume * s.am.masterVolume)

	player.Play()

	return nil
}

// LoadSound loads a sound effect from a file
func (am *AssetManager) LoadSound(name, filePath string) error {
	var format string
	ext := strings.ToLower(filepath.Ext(filePath)) // ".wav" or ".ogg"
	switch ext {
	case ".wav":
		format = "wav"
	case ".ogg":
		format = "ogg"
	default:
		return fmt.Errorf("unsupported audio format: %s", ext)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading sound file: %v", err)
	}

	return am.LoadSoundFromBytes(name, data, format)
}

// LoadSoundFromBytes loads a sound effect from byte data
func (am *AssetManager) LoadSoundFromBytes(name string, data []byte, format string) error {
	am.sounds[name] = &Sound{
		audioContext: am.audioContext,
		data:         data,
		format:       format,
		volume:       1.0,
		am:           am,
	}

	return nil
}

// GetSound returns a sound by name
func (am *AssetManager) GetSound(name string) *Sound {
	if sound, ok := am.sounds[name]; ok {
		return sound
	}
	return nil
}

// PlaySound plays a sound by name
func (am *AssetManager) PlaySound(name string) error {
	sound := am.GetSound(name)
	if sound == nil {
		return fmt.Errorf("sound not found: %s", name)
	}
	return sound.Play()
}
