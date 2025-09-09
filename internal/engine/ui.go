package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawScene(g.fb)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.scaleX, g.scaleY)
	screen.DrawImage(g.fb, op)

	g.drawHUD(screen)
	if g.minimap && g.state == statePlaying {
		g.drawMinimap(screen)
	}
	switch g.state {
	case stateMainMenu:
		g.drawMainMenu(screen)
	case stateInGameMenu:
		g.drawInGameMenu(screen)
	case stateOptions:
		g.drawOptionsMenu(screen)
	case stateStart:
		g.drawStart(screen)
	case stateMenu:
		g.drawMenu(screen)
	case stateLevelClear:
		g.drawLevelClear(screen)
	case stateGameOver:
		g.drawStateOverlay(screen, "YOU DIED", red)
	case stateWin:
		g.drawStateOverlay(screen, "YOU WIN!", uiAccent)
	}
}

func (g *Game) drawScene(dst *ebiten.Image) {
	for i := range g.zbuf {
		g.zbuf[i] = 1e9
	}
	g.drawFloorCeil(dst)
	g.drawWalls(dst)
	g.drawSprites(dst)
}

func (g *Game) drawHUD(dst *ebiten.Image) {
	if g.state == statePlaying {
		c := uiAccent
		h := 10
		w := 2
		cx := int(float64(renderW) / 2.0 * g.scaleX)
		cy := int(float64(renderH) / 2.0 * g.scaleY)
		drawRect(dst, g.pix, cx-w/2, cy-h, w, h*2, c)
		drawRect(dst, g.pix, cx-h, cy-w/2, h*2, w, c)
		if g.p.muzzleTime > 0 {
			a := uint8(80 * g.p.muzzleTime / 0.06)
			drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{255, 255, 200, a})
		}
	}

	gw, gh := 140, 60
	wx := (ScreenW - gw) / 2
	wy := ScreenH - gh - 8
	drawRect(dst, g.pix, wx, wy, gw, gh, color.RGBA{40, 40, 40, 230})
	drawRect(dst, g.pix, wx+gw/2-8, wy-20, 16, 24, color.RGBA{70, 70, 70, 230})

	// health/ammo
	barW := 220
	barH := 10
	bx := 12
	by := 16
	drawRect(dst, g.pix, bx-2, by-2, barW+4, barH+4, black)
	drawRect(dst, g.pix, bx, by, barW, barH, color.RGBA{60, 20, 20, 220})
	fill := int(float64(barW) * clamp01(float64(g.p.hp)/float64(playerMaxHP)))
	if fill > 0 {
		col := red
		if g.p.hp >= playerMaxHP/2 {
			col = green
		} else if g.p.hp > 20 {
			col = yellow
		}
		drawRect(dst, g.pix, bx, by, fill, barH, col)
	}
	text.Draw(dst, fmt.Sprintf("HP: %d / %d", g.p.hp, playerMaxHP), g.face, bx, by+barH+14, white)
	text.Draw(dst, fmt.Sprintf("Ammo: %d", g.p.ammo), g.face, bx, by+barH+30, yellow)

	// level & counters
	lx := ScreenW - 260
	ly := 20
	drawRect(dst, g.pix, lx-10, ly-16, 240, 56, color.RGBA{0, 0, 0, 160})
	text.Draw(dst, fmt.Sprintf("Level: %d / %d", g.level, g.totalLevels), g.face, lx, ly, uiAccent)
	ly += 18
	remaining := 0
	for _, e := range g.enemies {
		if !e.dead {
			remaining++
		}
	}
	text.Draw(dst, fmt.Sprintf("Defeated: %d", g.defeated), g.face, lx, ly, white)
	ly += 18
	text.Draw(dst, fmt.Sprintf("Remaining: %d", remaining), g.face, lx, ly, white)

	// Draw pickup messages
	g.drawPickupMessages(dst)
}

func (g *Game) drawPickupMessages(dst *ebiten.Image) {
	if len(g.pickupMessages) == 0 {
		return
	}

	// Position messages in the center-right area of the screen
	startX := ScreenW - 200
	startY := 100

	for i, msg := range g.pickupMessages {
		// Calculate alpha based on remaining time (fade out in last 0.5 seconds)
		alpha := uint8(255)
		if msg.timeLeft < 0.5 {
			alpha = uint8(255 * (msg.timeLeft / 0.5))
		}

		// Create color with alpha
		msgColor := color.RGBA{msg.color.R, msg.color.G, msg.color.B, alpha}

		// Draw message with slight offset for multiple messages
		y := startY + (i * 25)
		text.Draw(dst, msg.text, g.face, startX, y, msgColor)
	}
}

func (g *Game) drawMainMenu(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180})

	w, h := 400, 300
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	// Title
	lx := x + 18
	ly := y + 40
	text.Draw(dst, "DOOMLIKE", g.face, lx, ly, uiAccent)
	ly += 50

	// Menu options
	options := []string{"Start Game", "Options", "Quit"}
	for i, option := range options {
		color := white
		if i == g.menu.selectedOption {
			color = yellow
			// Draw selection indicator
			text.Draw(dst, ">", g.face, lx-20, ly, color)
		}
		text.Draw(dst, option, g.face, lx, ly, color)
		ly += 30
	}

	// Instructions
	ly += 20
	text.Draw(dst, "Use ↑/↓ to navigate, Enter to select", g.face, lx, ly, gray)
	ly += 20
	text.Draw(dst, "Esc to quit", g.face, lx, ly, gray)
}

func (g *Game) drawInGameMenu(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180})

	w, h := 400, 250
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	// Title
	lx := x + 18
	ly := y + 40
	text.Draw(dst, "PAUSED", g.face, lx, ly, uiAccent)
	ly += 50

	// Menu options
	options := []string{"Resume Game", "Options", "Quit Game"}
	for i, option := range options {
		color := white
		if i == g.menu.selectedInGameOption {
			color = yellow
			// Draw selection indicator
			text.Draw(dst, ">", g.face, lx-20, ly, color)
		}
		text.Draw(dst, option, g.face, lx, ly, color)
		ly += 30
	}

	// Instructions
	ly += 20
	text.Draw(dst, "Use ↑/↓ to navigate, Enter to select", g.face, lx, ly, gray)
	ly += 20
	text.Draw(dst, "Esc to quit game", g.face, lx, ly, gray)
}

func (g *Game) drawOptionsMenu(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180})

	w, h := 500, 350
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	// Title
	lx := x + 18
	ly := y + 40
	text.Draw(dst, "OPTIONS", g.face, lx, ly, uiAccent)
	ly += 50

	// Fire Rate Slider
	g.drawSlider(dst, lx, ly, 400, 50, minFireRate, maxFireRate, g.settings.fireRate, "Fire Rate:", g.menu.selectedSetting == 0)
	ly += 60

	// Other settings (non-slider) - conditionally show level count
	otherSettings := []struct {
		name  string
		value string
	}{
		{"Bullet Speed:", fmt.Sprintf("%.0f", g.settings.bulletSpeed)},
	}

	// Only show level count if we came from the main menu
	if g.previousState == stateMainMenu {
		otherSettings = append(otherSettings, struct {
			name  string
			value string
		}{"Level Count:", fmt.Sprintf("%d", g.settings.levelCount)})
	}

	for i, setting := range otherSettings {
		settingIndex := i + 1 // Offset by 1 since fire rate is index 0
		color := white
		if settingIndex == g.menu.selectedSetting {
			color = yellow
			// Draw selection indicator
			text.Draw(dst, ">", g.face, lx-20, ly, color)
		}
		text.Draw(dst, setting.name, g.face, lx, ly, color)
		text.Draw(dst, setting.value, g.face, lx+200, ly, color)
		ly += 35
	}

	// Instructions
	ly += 30
	text.Draw(dst, "Use ↑/↓ to navigate, ←/→ to adjust", g.face, lx, ly, gray)
	ly += 20
	text.Draw(dst, "Click on fire rate slider to set value", g.face, lx, ly, gray)
	ly += 20
	text.Draw(dst, "Esc to return to main menu", g.face, lx, ly, gray)
}

// drawSlider draws a slider with ticks at 10% increments
func (g *Game) drawSlider(dst *ebiten.Image, x, y, width, height int, minVal, maxVal, currentVal float64, label string, selected bool) {
	// Draw label
	labelColor := white
	if selected {
		labelColor = yellow
		// Draw selection indicator
		text.Draw(dst, ">", g.face, x-20, y, labelColor)
	}
	text.Draw(dst, label, g.face, x, y, labelColor)

	// Calculate slider position
	sliderY := y + 20
	sliderHeight := 8
	sliderWidth := width - 200 // Leave space for value display

	// Draw slider track background
	drawRect(dst, g.pix, x, sliderY, sliderWidth, sliderHeight, color.RGBA{60, 60, 60, 255})

	// Calculate current position (0.0 to 1.0) - reversed so left is fast, right is slow
	normalizedVal := (maxVal - currentVal) / (maxVal - minVal)
	if normalizedVal < 0 {
		normalizedVal = 0
	}
	if normalizedVal > 1 {
		normalizedVal = 1
	}

	// Draw slider fill (left side)
	fillWidth := int(float64(sliderWidth) * normalizedVal)
	if fillWidth > 0 {
		drawRect(dst, g.pix, x, sliderY, fillWidth, sliderHeight, color.RGBA{100, 150, 255, 255})
	}

	// Draw slider handle
	handleX := x + fillWidth - 4
	if handleX < x {
		handleX = x
	}
	if handleX > x+sliderWidth-8 {
		handleX = x + sliderWidth - 8
	}
	drawRect(dst, g.pix, handleX, sliderY-2, 8, sliderHeight+4, color.RGBA{200, 200, 200, 255})

	// Draw ticks at each 0.05s increment
	tickColor := color.RGBA{120, 120, 120, 255}
	for val := minVal; val <= maxVal; val += 0.05 {
		// Calculate position for this value (reversed)
		valNormalized := (maxVal - val) / (maxVal - minVal)
		if valNormalized >= 0 && valNormalized <= 1 {
			tickX := x + int(float64(sliderWidth)*valNormalized)
			tickY := sliderY + sliderHeight + 2
			drawRect(dst, g.pix, tickX, tickY, 1, 4, tickColor)
		}
	}

	// Draw value
	valueText := fmt.Sprintf("%.2fs", currentVal)
	text.Draw(dst, valueText, g.face, x+sliderWidth+10, y+15, white)
}
