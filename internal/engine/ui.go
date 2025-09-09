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
