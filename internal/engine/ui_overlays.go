package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) drawStart(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180})

	w, h := 620, 220
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	lx := x + 18
	ly := y + 46
	text.Draw(dst, "DOOMLIKE", g.face, lx, ly, uiAccent)
	ly += 26
	text.Draw(dst, fmt.Sprintf("Select number of levels:  %d", g.totalLevels), g.face, lx, ly, white)
	ly += 22
	text.Draw(dst, fmt.Sprintf("Min 1, Max %d", MaxLevelCap), g.face, lx, ly, gray)
	ly += 22
	text.Draw(dst, "←/→ or ↑/↓ to adjust, 0-9 to quick-set (0=10), Enter to start", g.face, lx, ly, yellow)
}

func (g *Game) drawMenu(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 128})

	w, h := 520, 200
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	lx := x + 18
	ly := y + 36
	text.Draw(dst, "Paused", g.face, lx, ly, uiAccent)
	ly += 30
	text.Draw(dst, "Enter: Resume", g.face, lx, ly, white)
	ly += 20
	text.Draw(dst, "Esc or Q: Quit Game", g.face, lx, ly, white)
	ly += 20
	text.Draw(dst, "WASD/Mouse | LMB/Space Shoot | M Minimap", g.face, lx, ly, white)
}

func (g *Game) drawStateOverlay(dst *ebiten.Image, title string, titleCol color.Color) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 160})

	w, h := 520, 200
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	lx := x + 18
	ly := y + 48
	text.Draw(dst, title, g.face, lx, ly, titleCol)
	ly += 36
	text.Draw(dst, "Press Enter to restart", g.face, lx, ly, white)
}

func (g *Game) drawLevelClear(dst *ebiten.Image) {
	drawRect(dst, g.pix, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180})

	w, h := 560, 220
	x := (ScreenW - w) / 2
	y := (ScreenH - h) / 2

	drawRect(dst, g.pix, x, y, w, h, uiBox)
	drawRect(dst, g.pix, x, y, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y+h-2, w, 2, uiAccent)
	drawRect(dst, g.pix, x, y, 2, h, uiAccent)
	drawRect(dst, g.pix, x+w-2, y, 2, h, uiAccent)

	lx := x + 18
	ly := y + 44
	title := fmt.Sprintf("LEVEL %d CLEARED!", g.level)
	text.Draw(dst, title, g.face, lx, ly, uiAccent)

	ly += 26
	text.Draw(dst, fmt.Sprintf("Defeated this run: %d", g.defeated), g.face, lx, ly, white)

	ly += 26
	if g.level < g.totalLevels {
		text.Draw(dst, fmt.Sprintf("Up next: Level %d / %d", g.level+1, g.totalLevels), g.face, lx, ly, white)
		ly += 22
		text.Draw(dst, "Press Enter to begin the next level", g.face, lx, ly, yellow)
	} else {
		text.Draw(dst, "Press Enter", g.face, lx, ly, yellow)
	}
}
