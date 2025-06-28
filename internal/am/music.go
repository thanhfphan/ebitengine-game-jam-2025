package am

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Music represents background music that can loop
type Music struct {
	audioContext *audio.Context
	data         []byte
	format       string
	volume       float64
	player       *audio.Player
	isPlaying    bool
	loop         bool
	am           *AssetManager
}

// Play starts playing the music
func (m *Music) Play() error {
	if m.isPlaying {
		return nil
	}

	if m.player == nil {
		if err := m.createPlayer(); err != nil {
			return err
		}
	}

	m.player.SetVolume(m.volume * m.am.masterVolume * m.am.musicVolume)

	m.player.Play()

	m.isPlaying = true
	return nil
}

// Stop stops the music playback
func (m *Music) Stop() error {
	if !m.isPlaying || m.player == nil {
		return nil
	}

	m.player.Pause()

	m.isPlaying = false
	return nil
}

// SetVolume sets the volume for this music track
func (m *Music) SetVolume(volume float64) {
	m.volume = volume
	if m.player != nil && m.isPlaying {
		m.player.SetVolume(m.volume * m.am.masterVolume * m.am.musicVolume)
	}
}

// IsPlaying returns whether the music is currently playing
func (m *Music) IsPlaying() bool {
	return m.isPlaying
}

// SetLoop sets whether the music should loop
func (m *Music) SetLoop(loop bool) {
	m.loop = loop
}

func (m *Music) Update() error {
	// Only needed for WAV (manual loop)
	if m.format == "wav" && m.loop && m.player != nil && !m.player.IsPlaying() {
		if err := m.createPlayer(); err != nil {
			return err
		}
		m.player.SetVolume(m.volume * m.am.masterVolume * m.am.musicVolume)
		m.player.Play()
	}
	return nil
}

// createPlayer creates a new audio player for the music
func (m *Music) createPlayer() error {
	var (
		err    error
		reader io.Reader
	)

	switch m.format {
	case "ogg":
		baseReader, err := vorbis.DecodeWithSampleRate(m.audioContext.SampleRate(), bytes.NewReader(m.data))
		if err != nil {
			return fmt.Errorf("error decoding OGG: %v", err)
		}

		if m.loop {
			reader = audio.NewInfiniteLoop(baseReader, baseReader.Length())
		} else {
			reader = baseReader
		}

	case "wav":
		reader, err = wav.DecodeWithSampleRate(m.audioContext.SampleRate(), bytes.NewReader(m.data))
		if err != nil {
			return fmt.Errorf("error decoding WAV: %v", err)
		}

	default:
		return fmt.Errorf("unsupported audio format: %s", m.format)
	}

	m.player, err = m.audioContext.NewPlayer(reader)
	if err != nil {
		return fmt.Errorf("error creating audio player: %v", err)
	}

	return nil
}

// LoadMusicFromBytes loads a music track from byte data
func (am *AssetManager) LoadMusicFromBytes(name string, data []byte, format string) error {
	am.music[name] = &Music{
		audioContext: am.audioContext,
		data:         data,
		format:       format,
		volume:       1.0,
		isPlaying:    false,
		loop:         true, // Default to looping for music
		am:           am,
	}

	return nil
}

// GetMusic returns a music track by name
func (am *AssetManager) GetMusic(name string) *Music {
	if music, ok := am.music[name]; ok {
		return music
	}
	return nil
}

// PlayMusic plays a music track by name
func (am *AssetManager) PlayMusic(name string) error {
	music := am.GetMusic(name)
	if music == nil {
		return fmt.Errorf("music not found: %s", name)
	}
	return music.Play()
}

// StopMusic stops a music track by name
func (am *AssetManager) StopMusic(name string) error {
	music := am.GetMusic(name)
	if music == nil {
		return fmt.Errorf("music not found: %s", name)
	}
	return music.Stop()
}

// StopAllMusic stops all currently playing music
func (am *AssetManager) StopAllMusic() {
	for _, music := range am.music {
		if music.IsPlaying() {
			music.Stop()
		}
	}
}
