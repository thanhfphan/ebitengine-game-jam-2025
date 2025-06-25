package am

import (
	"bytes"
	"fmt"
	"os"

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
		wavReader, err := wav.DecodeWithSampleRate(s.audioContext.SampleRate(), bytes.NewReader(s.data))
		if err != nil {
			return fmt.Errorf("error decoding WAV: %v", err)
		}
		player, err = s.audioContext.NewPlayer(wavReader)
	case "ogg":
		oggReader, err := vorbis.DecodeWithSampleRate(s.audioContext.SampleRate(), bytes.NewReader(s.data))
		if err != nil {
			return fmt.Errorf("error decoding OGG: %v", err)
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

	go func() {
		player.Play()
		for player.IsPlaying() {
			// Wait for playback to complete
		}
		player.Close()
	}()

	return nil
}

// LoadSound loads a sound effect from a file
func (am *AssetManager) LoadSound(name, filePath string) error {
	var format string
	if len(filePath) > 4 {
		ext := filePath[len(filePath)-3:]
		if ext == "wav" {
			format = "wav"
		} else if ext == "ogg" {
			format = "ogg"
		} else {
			return fmt.Errorf("unsupported audio format: %s", ext)
		}
	} else {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading sound file: %v", err)
	}

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
