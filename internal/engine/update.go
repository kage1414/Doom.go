package engine

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Layout(outW, outH int) (int, int) {
	if outW > 0 && outH > 0 {
		g.scaleX = float64(outW) / float64(renderW)
		g.scaleY = float64(outH) / float64(renderH)
	}
	return ScreenW, ScreenH
}

func (g *Game) Update() error {
	// Check for quit flag
	if g.shouldQuit {
		return ebiten.Termination
	}

	// Global Esc behavior
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch g.state {
		case statePlaying:
			g.state = stateInGameMenu
			g.mouseGrabbed = false
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		case stateInGameMenu:
			// Resume the game
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		case stateMenu:
			return ebiten.Termination
		case stateOptions:
			// Return to the previous state
			g.state = g.previousState
			g.menu.selectedSetting = 0
		}
	}

	switch g.state {
	case stateMainMenu:
		g.updateMainMenu()
		return nil

	case stateInGameMenu:
		g.updateInGameMenu()
		return nil

	case stateOptions:
		g.updateOptionsMenu()
		return nil

	case stateStart:
		// Choose total levels before starting
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if g.totalLevels < MaxLevelCap {
				g.totalLevels++
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if g.totalLevels > 1 {
				g.totalLevels--
			}
		}
		// Digit keys quick-set
		for key := ebiten.Key0; key <= ebiten.Key9; key++ {
			if inpututil.IsKeyJustPressed(key) {
				n := int(key - ebiten.Key0)
				if n == 0 {
					n = 10
				}
				if n > MaxLevelCap {
					n = MaxLevelCap
				}
				g.totalLevels = n
			}
		}
		// Enter to begin
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.level = 1
			g.setupLevel(g.level, true)
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		return nil

	case stateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return ebiten.Termination
		}
		return nil

	case stateLevelClear:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.level++
			if g.level > g.totalLevels {
				g.state = stateWin
				return nil
			}
			g.setupLevel(g.level, false)
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		return nil

	case stateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
		return nil

	case stateWin:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.resetToMainMenu()
		}
		return nil

	case statePlaying:
		dt := 1.0 / 60.0

		if g.p.cooldown > 0 {
			g.p.cooldown -= dt
			if g.p.cooldown < 0 {
				g.p.cooldown = 0
			}
		}
		if g.p.muzzleTime > 0 {
			g.p.muzzleTime -= dt
			if g.p.muzzleTime < 0 {
				g.p.muzzleTime = 0
			}
		}

		for _, e := range g.enemies {
			if e.blink > 0 {
				e.blink -= dt
				if e.blink < 0 {
					e.blink = 0
				}
			}
			if !e.dead {
				e.aiTime += dt
			}
		}

		g.updateProjectiles(dt)

		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			g.minimap = !g.minimap
		}

		if !g.mouseGrabbed {
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}

		if g.mouseGrabbed {
			x, _ := ebiten.CursorPosition()
			if g.lastMouseX != 0 {
				dx := x - g.lastMouseX
				g.p.angle += float64(dx) * mouseSens
				g.p.angle = normalizeAngle(g.p.angle)
			}
			g.lastMouseX = x
		} else {
			g.lastMouseX = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.p.angle -= rotSpeed * dt
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.p.angle += rotSpeed * dt
		}
		g.p.angle = normalizeAngle(g.p.angle)

		forward, side := 0.0, 0.0
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			forward++
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			forward--
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			side--
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			side++
		}
		speed := moveSpeed
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			speed *= sprintMul
		}
		if forward != 0 || side != 0 {
			l := math.Hypot(forward, side)
			forward /= l
			side /= l
			fx := math.Cos(g.p.angle)
			fy := math.Sin(g.p.angle)
			rx := -fy
			ry := fx
			vx := (fx*forward + rx*side) * speed * dt
			vy := (fy*forward + ry*side) * speed * dt
			g.moveWithCollision(vx, vy)
		}

		if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace)) &&
			g.p.cooldown <= 0 && g.p.ammo > 0 {
			g.p.cooldown = g.settings.fireRate
			g.p.muzzleTime = 0.06
			g.p.ammo--
			g.firePlayerShot()
		}

		for _, e := range g.enemies {
			if e.dead {
				continue
			}
			switch e.etype {
			case eZombie:
				g.seekEnemy(e, zombieSpeed, dt)
				if dist2(e.pos.x, e.pos.y, g.p.pos.x, g.p.pos.y) < (0.25+0.25)*(0.25+0.25) {
					g.p.hp -= int(touchDPS * dt)
					if g.p.hp < 0 {
						g.p.hp = 0
					}
				}
			case eRunner:
				g.seekEnemy(e, runnerSpeed, dt)
				if dist2(e.pos.x, e.pos.y, g.p.pos.x, g.p.pos.y) < (0.25+0.25)*(0.25+0.25) {
					g.p.hp -= int(touchDPS * dt)
					if g.p.hp < 0 {
						g.p.hp = 0
					}
				}
			case eShooter:
				g.shooterAI(e, dt)
			}
		}

		for _, pk := range g.pickups {
			if pk.took {
				continue
			}
			if dist2(pk.pos.x, pk.pos.y, g.p.pos.x, g.p.pos.y) < 0.5*0.5 {
				switch pk.ptype {
				case pickupMedkit:
					if g.p.hp < playerMaxHP {
						g.p.hp += medkitHeal
						if g.p.hp > playerMaxHP {
							g.p.hp = playerMaxHP
						}
						pk.took = true
						// Add pickup message
						g.pickupMessages = append(g.pickupMessages, pickupMessage{
							text:     fmt.Sprintf("+%d Health", medkitHeal),
							color:    green,
							timeLeft: pickupMessageDuration,
						})
					}
				case pickupAmmo:
					g.p.ammo += ammoPickupAmt
					pk.took = true
					// Add pickup message
					g.pickupMessages = append(g.pickupMessages, pickupMessage{
						text:     fmt.Sprintf("+%d Ammo", ammoPickupAmt),
						color:    yellow,
						timeLeft: pickupMessageDuration,
					})
				}
			}
		}

		if g.p.hp <= 0 {
			g.state = stateGameOver
			g.mouseGrabbed = false
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			return nil
		}

		allDead := true
		for _, e := range g.enemies {
			if !e.dead {
				allDead = false
				break
			}
		}
		if allDead {
			g.advanceLevelOrWin()
			return nil
		}

		// Update pickup messages
		g.updatePickupMessages(dt)
	}
	return nil
}

func (g *Game) moveWithCollision(dx, dy float64) {
	newX := g.p.pos.x + dx
	newY := g.p.pos.y + dy
	if !g.isSolid(int(math.Floor(newX)), int(math.Floor(g.p.pos.y))) {
		g.p.pos.x = newX
	}
	if !g.isSolid(int(math.Floor(g.p.pos.x)), int(math.Floor(newY))) {
		g.p.pos.y = newY
	}
}

func (g *Game) updatePickupMessages(dt float64) {
	// Update message timers and remove expired messages
	nm := g.pickupMessages[:0]
	for _, msg := range g.pickupMessages {
		msg.timeLeft -= dt
		if msg.timeLeft > 0 {
			nm = append(nm, msg)
		}
	}
	g.pickupMessages = nm
}

func (g *Game) updateMainMenu() {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Navigate menu options with keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.menu.selectedOption--
		if g.menu.selectedOption < 0 {
			g.menu.selectedOption = 2 // Wrap to last option
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.menu.selectedOption++
		if g.menu.selectedOption > 2 {
			g.menu.selectedOption = 0 // Wrap to first option
		}
	}

	// Handle mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		clickedOption := g.getMainMenuOptionAt(g.mouseX, g.mouseY)
		if clickedOption >= 0 {
			g.menu.selectedOption = clickedOption
			// Trigger the selection immediately
			g.selectMainMenuOption()
		}
	}

	// Select option with Enter key
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.selectMainMenuOption()
	}
}

func (g *Game) updateInGameMenu() {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Navigate menu options with keyboard
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.menu.selectedInGameOption--
		if g.menu.selectedInGameOption < 0 {
			g.menu.selectedInGameOption = 2 // Wrap to last option
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.menu.selectedInGameOption++
		if g.menu.selectedInGameOption > 2 {
			g.menu.selectedInGameOption = 0 // Wrap to first option
		}
	}

	// Handle mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		clickedOption := g.getInGameMenuOptionAt(g.mouseX, g.mouseY)
		if clickedOption >= 0 {
			g.menu.selectedInGameOption = clickedOption
			// Trigger the selection immediately
			g.selectInGameMenuOption()
		}
	}

	// Select option with Enter key
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.selectInGameMenuOption()
	}
}

func (g *Game) updateOptionsMenu() {
	// Update mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Calculate max setting index based on context
	maxSetting := 1 // Fire rate (0) + Bullet speed (1)
	if g.previousState == stateMainMenu {
		maxSetting = 2 // Fire rate (0) + Bullet speed (1) + Level count (2)
	}

	// Ensure selected setting is valid for current context
	if g.menu.selectedSetting > maxSetting {
		g.menu.selectedSetting = maxSetting
	}

	// Navigate settings
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.menu.selectedSetting--
		if g.menu.selectedSetting < 0 {
			g.menu.selectedSetting = maxSetting // Wrap to last setting
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.menu.selectedSetting++
		if g.menu.selectedSetting > maxSetting {
			g.menu.selectedSetting = 0 // Wrap to first setting
		}
	}

	// Handle mouse clicks on fire rate slider
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.menu.selectedSetting == 0 {
		g.handleSliderClick()
	}

	// Adjust settings
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		delta := 1.0
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			delta = -1.0
		}

		switch g.menu.selectedSetting {
		case 0: // Fire Rate (reversed: left = faster, right = slower)
			g.settings.fireRate -= delta * 0.05
			if g.settings.fireRate < minFireRate {
				g.settings.fireRate = minFireRate
			}
			if g.settings.fireRate > maxFireRate {
				g.settings.fireRate = maxFireRate
			}
			g.saveSettings()
		case 1: // Bullet Speed
			g.settings.bulletSpeed += delta * 2.0
			if g.settings.bulletSpeed < minBulletSpeed {
				g.settings.bulletSpeed = minBulletSpeed
			}
			if g.settings.bulletSpeed > maxBulletSpeed {
				g.settings.bulletSpeed = maxBulletSpeed
			}
			g.saveSettings()
		case 2: // Level Count (only available from main menu)
			if g.previousState == stateMainMenu {
				g.settings.levelCount += int(delta)
				if g.settings.levelCount < minLevelCount {
					g.settings.levelCount = minLevelCount
				}
				if g.settings.levelCount > maxLevelCount {
					g.settings.levelCount = maxLevelCount
				}
				g.saveSettings()
			}
		}
	}
}

// saveSettings saves the current settings to the database
func (g *Game) saveSettings() {
	if g.db != nil {
		if err := g.db.SaveSettings(&g.settings); err != nil {
			log.Printf("Failed to save settings: %v", err)
		}
	}
}

// handleSliderClick handles mouse clicks on the fire rate slider
func (g *Game) handleSliderClick() {
	// Calculate slider bounds (matching the drawSlider function)
	// These coordinates should match the slider position in drawOptionsMenu
	sliderX := 18 + 18       // lx from drawOptionsMenu
	sliderY := 40 + 50 + 20  // ly + 20 from drawOptionsMenu
	sliderWidth := 400 - 200 // width - 200 from drawSlider

	// Check if click is within slider bounds
	if g.mouseX >= sliderX && g.mouseX <= sliderX+sliderWidth &&
		g.mouseY >= sliderY && g.mouseY <= sliderY+8 {

		// Calculate new value based on click position (center-out)
		clickPos := float64(g.mouseX - sliderX)
		normalizedPos := clickPos / float64(sliderWidth)
		if normalizedPos < 0 {
			normalizedPos = 0
		}
		if normalizedPos > 1 {
			normalizedPos = 1
		}

		// Convert to fire rate value (reversed: left=fast, right=slow)
		newFireRate := maxFireRate - normalizedPos*(maxFireRate-minFireRate)
		g.settings.fireRate = newFireRate
		g.saveSettings()
	}
}

// selectMainMenuOption handles the main menu option selection
func (g *Game) selectMainMenuOption() {
	switch g.menu.selectedOption {
	case 0: // Start Game
		g.totalLevels = g.settings.levelCount
		g.level = 1
		g.setupLevel(g.level, true)
		g.state = statePlaying
		g.mouseGrabbed = true
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		g.lastMouseX = 0
	case 1: // Options
		g.previousState = stateMainMenu
		g.state = stateOptions
		g.menu.selectedSetting = 0
	case 2: // Quit
		g.shouldQuit = true
	}
}

// getMainMenuOptionAt returns the menu option index at the given mouse coordinates, or -1 if none
func (g *Game) getMainMenuOptionAt(mouseX, mouseY int) int {
	// Menu dimensions and position (same as in drawMainMenu)
	w, h := 400, 300
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	// Check if click is within menu bounds
	if mouseX < x || mouseX > x+w || mouseY < y || mouseY > y+h {
		return -1
	}

	// Calculate option positions
	ly := y + 40 + 50 // Start after title

	// Check each option (3 options, 30 pixels apart)
	for i := 0; i < 3; i++ {
		optionY := ly + i*30
		if mouseY >= optionY-15 && mouseY <= optionY+15 {
			return i
		}
	}

	return -1
}

// selectInGameMenuOption handles the in-game menu option selection
func (g *Game) selectInGameMenuOption() {
	switch g.menu.selectedInGameOption {
	case 0: // Resume Game
		g.state = statePlaying
		g.mouseGrabbed = true
		ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		g.lastMouseX = 0
	case 1: // Options
		g.previousState = stateInGameMenu
		g.state = stateOptions
		g.menu.selectedSetting = 0 // Reset to fire rate (valid for both contexts)
	case 2: // Quit Game
		g.resetToMainMenu()
	}
}

// getInGameMenuOptionAt returns the in-game menu option index at the given mouse coordinates, or -1 if none
func (g *Game) getInGameMenuOptionAt(mouseX, mouseY int) int {
	// Menu dimensions and position (same as in drawInGameMenu)
	w, h := 400, 250
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	// Check if click is within menu bounds
	if mouseX < x || mouseX > x+w || mouseY < y || mouseY > y+h {
		return -1
	}

	// Calculate option positions
	ly := y + 40 + 50 // Start after title

	// Check each option (3 options, 30 pixels apart)
	for i := 0; i < 3; i++ {
		optionY := ly + i*30
		if mouseY >= optionY-15 && mouseY <= optionY+15 {
			return i
		}
	}

	return -1
}
