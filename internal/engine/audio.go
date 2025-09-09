package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// generateGunshotSound creates a simple gunshot sound programmatically
func generateGunshotSound(sampleRate int) []byte {
	const duration = 0.05 // 50ms - shorter for more responsive feel
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Create a gunshot-like sound with noise and decay
		noise := (math.Sin(t*2000) + math.Sin(t*3000) + math.Sin(t*4000)) / 3
		decay := math.Exp(-t * 20) // Exponential decay
		envelope := noise * decay

		// Convert to 16-bit PCM
		sample := int16(envelope * 16000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// initAudio initializes the audio system
func (g *Game) initAudio() error {
	// Create audio context
	g.audioContext = audio.NewContext(44100)

	// Generate gunshot sound data
	g.bulletSoundData = generateGunshotSound(44100)

	return nil
}

// playBulletSound plays the bullet firing sound
func (g *Game) playBulletSound() {
	if g.audioContext != nil && g.bulletSoundData != nil {
		// Create a new player for each shot to avoid rewind delay
		player := audio.NewPlayerFromBytes(g.audioContext, g.bulletSoundData)
		player.SetVolume(0.3)
		player.Play()
	}
}
