package engine

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func deg2rad(d float64) float64 { return d * math.Pi / 180.0 }

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func shade(c color.RGBA, mul float64) color.RGBA {
	r := uint8(clamp01(float64(c.R)*mul/255.0) * 255.0)
	g := uint8(clamp01(float64(c.G)*mul/255.0) * 255.0)
	b := uint8(clamp01(float64(c.B)*mul/255.0) * 255.0)
	return color.RGBA{r, g, b, c.A}
}

func drawRect(dst, pix *ebiten.Image, x, y, w, h int, clr color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(h))
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorM.ScaleWithColor(clr)
	dst.DrawImage(pix, op)
}

func dist2(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

func mapSideVertical(stepX int) int {
	if stepX > 0 {
		return 0
	}
	return 1
}

func mapSideHorizontal(stepY int) int {
	if stepY > 0 {
		return 2
	}
	return 3
}

func normalizeAngle(a float64) float64 {
	for a < 0 {
		a += 2 * math.Pi
	}
	for a >= 2*math.Pi {
		a -= 2 * math.Pi
	}
	return a
}

func hWorldCellAtRay(g *Game, angle float64, dist float64) (int, int) {
	wx := g.p.pos.x + math.Cos(angle)*dist
	wy := g.p.pos.y + math.Sin(angle)*dist
	return int(math.Floor(wx)), int(math.Floor(wy))
}
