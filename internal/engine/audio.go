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

// generateCoinSound creates a satisfying coin/ding sound for enemy kills
func generateCoinSound(sampleRate int) []byte {
	const duration = 0.3 // 300ms
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Create a bright, bell-like sound with multiple harmonics
		// Main frequency starts high and drops slightly (like a coin spinning)
		baseFreq := 800.0 - 200.0*t // 800Hz to 600Hz

		// Add harmonics for richness
		fundamental := math.Sin(t * baseFreq * 2 * math.Pi)
		harmonic2 := 0.5 * math.Sin(t*baseFreq*2*2*math.Pi)
		harmonic3 := 0.25 * math.Sin(t*baseFreq*3*2*math.Pi)
		harmonic4 := 0.125 * math.Sin(t*baseFreq*4*2*math.Pi)

		// Combine harmonics
		noise := fundamental + harmonic2 + harmonic3 + harmonic4

		// Add a slight metallic "ping" with higher frequency
		ping := 0.3 * math.Sin(t*1200*2*math.Pi) * math.Exp(-t*8)
		noise += ping

		// Envelope: quick attack, sustained, quick decay
		envelope := 1.0
		if t < 0.02 {
			envelope = t / 0.02 // Quick attack
		} else if t > 0.2 {
			envelope = (0.3 - t) / 0.1 // Quick decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 12000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateReloadSound creates a mechanical reload sound for ammo pickups
func generateReloadSound(sampleRate int) []byte {
	const duration = 0.4 // 400ms
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Create mechanical reload sound with multiple components
		// Magazine insertion sound (low frequency thunk)
		thunkFreq := 120.0 + 30.0*math.Sin(t*8) // 120-150 Hz
		thunk := math.Sin(t*thunkFreq*2*math.Pi) * math.Exp(-t*3)

		// Magazine click (higher frequency click)
		clickFreq := 400.0 + 100.0*math.Sin(t*15) // 400-500 Hz
		click := math.Sin(t*clickFreq*2*math.Pi) * (1.0 + 0.5*math.Sin(t*25))

		// Slide/bolt action (mid frequency mechanical sound)
		slideFreq := 250.0 + 50.0*math.Sin(t*6) // 250-300 Hz
		slide := math.Sin(t*slideFreq*2*math.Pi) * math.Sin(t*12)

		// Combine components
		noise := 0.4*thunk + 0.3*click + 0.3*slide

		// Envelope: quick attack, sustained, quick decay
		envelope := 1.0
		if t < 0.05 {
			envelope = t / 0.05 // Quick attack
		} else if t > 0.3 {
			envelope = (0.4 - t) / 0.1 // Quick decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 10000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateBulletWhizSound creates a whooshing sound for bullets passing by
func generateBulletWhizSound(sampleRate int) []byte {
	const duration = 0.15 // 150ms - short and sharp
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Create a whooshing sound with noise and frequency sweep
		// Start high frequency and sweep down (like a bullet whizzing by)
		baseFreq := 800.0 - 400.0*t // 800Hz to 400Hz sweep

		// Add white noise for the "whoosh" effect
		noise := (math.Sin(t*baseFreq*2*math.Pi) + math.Sin(t*baseFreq*1.5*2*math.Pi)) / 2

		// Add higher frequency components for the "whiz"
		whizFreq := 1200.0 - 600.0*t // 1200Hz to 600Hz
		whiz := 0.3 * math.Sin(t*whizFreq*2*math.Pi)

		// Combine noise and whiz
		sound := 0.7*noise + 0.3*whiz

		// Envelope: quick attack, sustained, quick decay
		envelope := 1.0
		if t < 0.02 {
			envelope = t / 0.02 // Quick attack
		} else if t > 0.1 {
			envelope = (0.15 - t) / 0.05 // Quick decay
		}

		// Convert to 16-bit PCM
		sample := int16(sound * envelope * 8000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateOneUpSound creates a classic 1-up sound for health pickups
func generateOneUpSound(sampleRate int) []byte {
	const duration = 0.8 // 800ms
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Create ascending 1-up sound (like Mario)
		// Start low and rise up in pitch
		baseFreq := 200.0 + 400.0*t // 200Hz to 600Hz

		// Add harmonics for richness
		fundamental := math.Sin(t * baseFreq * 2 * math.Pi)
		harmonic2 := 0.5 * math.Sin(t*baseFreq*2*2*math.Pi)
		harmonic3 := 0.25 * math.Sin(t*baseFreq*3*2*math.Pi)

		// Combine harmonics
		noise := fundamental + harmonic2 + harmonic3

		// Add a slight vibrato for musical quality
		vibrato := 0.1 * math.Sin(t*8*2*math.Pi)
		noise += vibrato

		// Envelope: slow attack, sustained, slow decay
		envelope := 1.0
		if t < 0.1 {
			envelope = t / 0.1 // Slow attack
		} else if t > 0.6 {
			envelope = (0.8 - t) / 0.2 // Slow decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 8000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateZombieGrumbler creates a low, guttural grumbling sound
func generateZombieGrumbler(sampleRate int) []byte {
	const duration = 2.0 // 2 seconds
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Low frequency grumbling with some variation
		baseFreq := 80.0 + 20.0*math.Sin(t*0.5) // 80-100 Hz base
		noise := math.Sin(t*baseFreq) + 0.3*math.Sin(t*baseFreq*2) + 0.1*math.Sin(t*baseFreq*3)

		// Add some random variation for guttural effect
		variation := 0.2 * math.Sin(t*15) * math.Sin(t*23)
		noise += variation

		// Envelope with slow attack and decay
		envelope := 1.0
		if t < 0.1 {
			envelope = t / 0.1 // Attack
		} else if t > 1.8 {
			envelope = (2.0 - t) / 0.2 // Decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 8000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateRunnerGrumbler creates a higher-pitched, faster grumbling sound
func generateRunnerGrumbler(sampleRate int) []byte {
	const duration = 1.5 // 1.5 seconds
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Higher frequency, more erratic
		baseFreq := 150.0 + 50.0*math.Sin(t*2.0) // 150-200 Hz base
		noise := math.Sin(t*baseFreq) + 0.4*math.Sin(t*baseFreq*1.5) + 0.2*math.Sin(t*baseFreq*2.5)

		// Add rapid variations for erratic effect
		variation := 0.3 * math.Sin(t*25) * math.Sin(t*37)
		noise += variation

		// Quick attack, sustained, quick decay
		envelope := 1.0
		if t < 0.05 {
			envelope = t / 0.05 // Quick attack
		} else if t > 1.3 {
			envelope = (1.5 - t) / 0.2 // Quick decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 6000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// generateShooterGrumbler creates a mechanical, clicking grumbling sound
func generateShooterGrumbler(sampleRate int) []byte {
	const duration = 2.5 // 2.5 seconds
	samples := int(float64(sampleRate) * duration)
	data := make([]byte, samples*2) // 16-bit audio

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)

		// Mechanical clicking with low rumble
		clickFreq := 200.0 + 100.0*math.Sin(t*0.8) // 200-300 Hz
		rumbleFreq := 60.0 + 20.0*math.Sin(t*0.3)  // 60-80 Hz rumble

		// Combine click and rumble
		click := math.Sin(t*clickFreq) * (1.0 + 0.5*math.Sin(t*50)) // Clicking pattern
		rumble := 0.6 * math.Sin(t*rumbleFreq)
		noise := click + rumble

		// Add mechanical variations
		variation := 0.2 * math.Sin(t*30) * math.Sin(t*47)
		noise += variation

		// Slow attack, sustained, slow decay
		envelope := 1.0
		if t < 0.2 {
			envelope = t / 0.2 // Slow attack
		} else if t > 2.0 {
			envelope = (2.5 - t) / 0.5 // Slow decay
		}

		// Convert to 16-bit PCM
		sample := int16(noise * envelope * 7000)
		data[i*2] = byte(sample)
		data[i*2+1] = byte(sample >> 8)
	}

	return data
}

// initAudio initializes the audio system
func (g *Game) initAudio() error {
	// Create audio context
	g.audioContext = audio.NewContext(44100)

	// Generate sound data
	g.bulletSoundData = generateGunshotSound(44100)
	g.coinSoundData = generateCoinSound(44100)
	g.reloadSoundData = generateReloadSound(44100)
	g.oneUpSoundData = generateOneUpSound(44100)
	g.bulletWhizData = generateBulletWhizSound(44100)

	// Create grumbling sound players
	zombieData := generateZombieGrumbler(44100)
	runnerData := generateRunnerGrumbler(44100)
	shooterData := generateShooterGrumbler(44100)

	g.zombieGrumbler = audio.NewPlayerFromBytes(g.audioContext, zombieData)
	g.runnerGrumbler = audio.NewPlayerFromBytes(g.audioContext, runnerData)
	g.shooterGrumbler = audio.NewPlayerFromBytes(g.audioContext, shooterData)

	// Set initial volumes (will be adjusted based on distance)
	g.zombieGrumbler.SetVolume(0.0)
	g.runnerGrumbler.SetVolume(0.0)
	g.shooterGrumbler.SetVolume(0.0)

	// Start all grumblers in loop mode
	g.zombieGrumbler.SetBufferSize(1024)
	g.runnerGrumbler.SetBufferSize(1024)
	g.shooterGrumbler.SetBufferSize(1024)

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

// playCoinSound plays the coin/ding sound for enemy kills
func (g *Game) playCoinSound() {
	if g.audioContext != nil && g.coinSoundData != nil {
		// Create a new player for each coin sound
		player := audio.NewPlayerFromBytes(g.audioContext, g.coinSoundData)
		player.SetVolume(0.1)
		player.Play()
	}
}

// playReloadSound plays the reload sound for ammo pickups
func (g *Game) playReloadSound() {
	if g.audioContext != nil && g.reloadSoundData != nil {
		// Create a new player for each reload sound
		player := audio.NewPlayerFromBytes(g.audioContext, g.reloadSoundData)
		player.SetVolume(0.1)
		player.Play()
	}
}

// playOneUpSound plays the 1-up sound for health pickups
func (g *Game) playOneUpSound() {
	if g.audioContext != nil && g.oneUpSoundData != nil {
		// Create a new player for each 1-up sound
		player := audio.NewPlayerFromBytes(g.audioContext, g.oneUpSoundData)
		player.SetVolume(0.1)
		player.Play()
	}
}

// playBulletWhizSound plays the bullet whiz sound for bullets passing by
func (g *Game) playBulletWhizSound() {
	if g.audioContext != nil && g.bulletWhizData != nil {
		// Create a new player for each whiz sound
		player := audio.NewPlayerFromBytes(g.audioContext, g.bulletWhizData)
		player.SetVolume(0.15)
		player.Play()
	}
}

// updateGrumblingSounds updates the volume of grumbling sounds based on distance to enemies
func (g *Game) updateGrumblingSounds() {
	if g.audioContext == nil {
		return
	}

	// Count enemies by type
	zombieCount := 0
	runnerCount := 0
	shooterCount := 0
	minZombieDist := 1000.0
	minRunnerDist := 1000.0
	minShooterDist := 1000.0

	// Calculate distances to nearest enemies of each type
	for _, enemy := range g.enemies {
		if enemy.dead {
			continue
		}

		dx := enemy.pos.x - g.p.pos.x
		dy := enemy.pos.y - g.p.pos.y
		dist := math.Sqrt(dx*dx + dy*dy)

		switch enemy.etype {
		case eZombie:
			zombieCount++
			if dist < minZombieDist {
				minZombieDist = dist
			}
		case eRunner:
			runnerCount++
			if dist < minRunnerDist {
				minRunnerDist = dist
			}
		case eShooter:
			shooterCount++
			if dist < minShooterDist {
				minShooterDist = dist
			}
		}
	}

	// Update zombie grumbling
	if zombieCount > 0 && g.zombieGrumbler != nil {
		volume := g.calculateGrumblingVolume(minZombieDist)
		g.zombieGrumbler.SetVolume(volume)
		if !g.zombieGrumbler.IsPlaying() {
			g.zombieGrumbler.Rewind()
			g.zombieGrumbler.Play()
		}
	} else if g.zombieGrumbler != nil {
		g.zombieGrumbler.SetVolume(0.0)
	}

	// Update runner grumbling
	if runnerCount > 0 && g.runnerGrumbler != nil {
		volume := g.calculateGrumblingVolume(minRunnerDist)
		g.runnerGrumbler.SetVolume(volume)
		if !g.runnerGrumbler.IsPlaying() {
			g.runnerGrumbler.Rewind()
			g.runnerGrumbler.Play()
		}
	} else if g.runnerGrumbler != nil {
		g.runnerGrumbler.SetVolume(0.0)
	}

	// Update shooter grumbling
	if shooterCount > 0 && g.shooterGrumbler != nil {
		volume := g.calculateGrumblingVolume(minShooterDist)
		g.shooterGrumbler.SetVolume(volume)
		if !g.shooterGrumbler.IsPlaying() {
			g.shooterGrumbler.Rewind()
			g.shooterGrumbler.Play()
		}
	} else if g.shooterGrumbler != nil {
		g.shooterGrumbler.SetVolume(0.0)
	}
}

// calculateGrumblingVolume calculates volume based on distance (closer = louder)
func (g *Game) calculateGrumblingVolume(distance float64) float64 {
	const maxDistance = 8.0 // Maximum distance to hear grumbling
	const maxVolume = 0.4   // Maximum volume when very close

	if distance > maxDistance {
		return 0.0
	}

	// Inverse square law for realistic distance falloff
	normalizedDist := distance / maxDistance
	volume := maxVolume * (1.0 - normalizedDist*normalizedDist)

	if volume < 0.0 {
		volume = 0.0
	}

	return volume
}
