package engine

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) drawMinimap(dst *ebiten.Image) {
	if g.mapW == 0 || g.mapH == 0 {
		return
	}

	// Calculate minimap size as 15% of screen to avoid HUD overlap
	screenW := dst.Bounds().Dx()
	minimapSize := int(float64(screenW) * 0.15) // 15% of screen width

	// Calculate scale to fit the map in the minimap size
	scale := minimapSize / g.mapW
	if scale < 1 {
		scale = 1
	}

	// Position in top-left corner with margin, avoiding HUD area
	margin := 20
	px := margin
	py := margin + 60 // Move down to avoid health/ammo display

	// frame
	drawRect(dst, g.pix, px-1, py-1, g.mapW*scale+2, g.mapH*scale+2, uiBox)

	// tiles
	for y := 0; y < g.mapH; y++ {
		for x := 0; x < g.mapW; x++ {
			idx := y*g.mapW + x
			if g.reachable != nil && !g.reachable[idx] {
				drawRect(dst, g.pix, px+x*scale, py+y*scale, scale, scale, color.RGBA{0, 0, 0, 255})
				continue
			}
			t := g.world[idx]
			col := color.RGBA{18, 50, 18, 255}
			if t == tWall {
				col = color.RGBA{120, 120, 120, 255}
			}
			drawRect(dst, g.pix, px+x*scale, py+y*scale, scale, scale, col)
		}
	}

	// player
	cx := px + int(g.p.pos.x*float64(scale))
	cy := py + int(g.p.pos.y*float64(scale))
	drawRect(dst, g.pix, cx-2, cy-2, 4, 4, uiAccent)

	// aim direction (short pointer; stops at walls)
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
		ix := int(math.Floor(wx))
		iy := int(math.Floor(wy))
		if g.isSolid(ix, iy) {
			break
		}
		pixX := px + int(wx*float64(scale))
		pixY := py + int(wy*float64(scale))
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
		ex := px + int(e.pos.x*float64(scale))
		ey := py + int(e.pos.y*float64(scale))
		drawRect(dst, g.pix, ex-2, ey-2, 4, 4, ec)
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
		pxx := px + int(pk.pos.x*float64(scale))
		pyy := py + int(pk.pos.y*float64(scale))
		drawRect(dst, g.pix, pxx-1, pyy-1, 2, 2, pc)
	}
}
