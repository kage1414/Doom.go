package engine

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

func (g *Game) drawMinimap(dst *ebiten.Image) {
	const scale = 3
	const px, py = 10, 80

	drawRect(dst, g.pix, px-1, py-1, g.mapW*scale+2, g.mapH*scale+2, uiBox)

	for y := 0; y < g.mapH; y++ {
		for x := 0; x < g.mapW; x++ {
			idx := y*g.mapW + x
			if g.reachable != nil && !g.reachable[idx] {
				drawRect(dst, g.pix, px+x*scale, py+y*scale, scale, scale, color.RGBA{0, 0, 0, 255})
				continue
			}
			t := g.world[idx]
			col := color.RGBA{20, 80, 20, 255}
			if t == tWall {
				col = color.RGBA{120, 120, 120, 255}
			}
			drawRect(dst, g.pix, px+x*scale, py+y*scale, scale, scale, col)
		}
	}

	// player position
	cx := px + int(g.p.pos.x*scale)
	cy := py + int(g.p.pos.y*scale)
	drawRect(dst, g.pix, cx-2, cy-2, 4, 4, uiAccent)

	// aim direction arrow (short)
	dirX := math.Cos(g.p.angle)
	dirY := math.Sin(g.p.angle)
	const arrowLen = 3.0 // tiles
	steps := int(arrowLen * 6)
	fx := g.p.pos.x
	fy := g.p.pos.y
	for i := 1; i <= steps; i++ {
		t := float64(i) / 6.0
		wx := fx + dirX*t
		wy := fy + dirY*t
		pixX := px + int(wx*float64(scale))
		pixY := py + int(wy*float64(scale))
		if g.isSolid(int(math.Floor(wx)), int(math.Floor(wy))) {
			break
		}
		drawRect(dst, g.pix, pixX-1, pixY-1, 2, 2, uiAccent)
	}

	// enemies
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		ec := gray
		switch e.etype {
		case eZombie:
			ec = gray
		case eRunner:
			ec = cyan
		case eShooter:
			ec = magenta
		}
		drawRect(dst, g.pix, px+int(e.pos.x*scale)-2, py+int(e.pos.y*scale)-2, 4, 4, ec)
	}

	// pickups
	for _, pk := range g.pickups {
		if pk.took {
			continue
		}
		pc := green
		if pk.ptype == pickupAmmo {
			pc = yellow
		}
		drawRect(dst, g.pix, px+int(pk.pos.x*scale)-1, py+int(pk.pos.y*scale)-1, 2, 2, pc)
	}
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

// New: level-complete overlay that pauses between levels until Enter is pressed
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
	if g.level < LevelMax {
		text.Draw(dst, fmt.Sprintf("Up next: Level %d", g.level+1), g.face, lx, ly, white)
		ly += 22
		text.Draw(dst, "Press Enter to begin the next level", g.face, lx, ly, yellow)
	} else {
		text.Draw(dst, "Press Enter", g.face, lx, ly, yellow)
	}
}
