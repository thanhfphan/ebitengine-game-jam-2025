package am

// SetMasterVolume sets the master volume for all audio
func (am *AssetManager) SetMasterVolume(volume float64) {
	volume = clamp(volume, 0, 1) // Clamp volume between 0 and 1
	am.masterVolume = volume

	for _, music := range am.music {
		if music.IsPlaying() {
			music.SetVolume(music.volume)
		}
	}
}

// GetMasterVolume returns the current master volume
func (am *AssetManager) GetMasterVolume() float64 {
	return am.masterVolume
}

// SetMusicVolume sets the volume for all music tracks
func (am *AssetManager) SetMusicVolume(volume float64) {
	volume = clamp(volume, 0, 1) // Clamp
	am.musicVolume = volume

	for _, music := range am.music {
		if music.IsPlaying() {
			music.SetVolume(music.volume)
		}
	}
}

// GetMusicVolume returns the current music volume
func (am *AssetManager) GetMusicVolume() float64 {
	return am.musicVolume
}

// Update should be called every frame to handle music looping
func (am *AssetManager) Update() {
	for _, music := range am.music {
		if music.IsPlaying() {
			music.Update()
		}
	}
}

func clamp(x, min, max float64) float64 {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
