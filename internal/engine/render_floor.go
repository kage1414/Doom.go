package engine

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) drawFloorCeil(dst *ebiten.Image) {
	h := float64(renderH)
	half := h / 2.0
	fov := deg2rad(fovDegrees)
	planeLen := math.Tan(fov / 2.0)
	dirX := math.Cos(g.p.angle)
	dirY := math.Sin(g.p.angle)
	planeX := -dirY * planeLen
	planeY := dirX * planeLen

	for sy := 0; sy < renderH; sy++ {
		row := float64(sy) - half
		if row == 0 {
			continue
		}
		invert := 1.0
		if row < 0 {
			invert = -1.0
		}
		rowDist := half / (invert * row)

		rayDirLX := dirX - planeX
		rayDirLY := dirY - planeY
		rayDirRX := dirX + planeX
		rayDirRY := dirY + planeY

		stepX := (rayDirRX - rayDirLX) * rowDist / float64(renderW)
		stepY := (rayDirRY - rayDirLY) * rowDist / float64(renderW)

		floorX := g.p.pos.x + rayDirLX*rowDist
		floorY := g.p.pos.y + rayDirLY*rowDist

		for sx := 0; sx < renderW; sx++ {
			wx := floorX
			wy := floorY

			cellX := int(math.Floor(wx))
			cellY := int(math.Floor(wy))

			// If world cell is out of bounds or NOT reachable from player, don't render it.
			if cellX < 0 || cellY < 0 || cellX >= g.mapW || cellY >= g.mapH ||
				!g.reachable[cellY*g.mapW+cellX] {
				dst.Set(sx, sy, black)
				floorX += stepX
				floorY += stepY
				continue
			}

			parity := (cellX+cellY)&1 == 0
			baseF := floorA
			if !parity {
				baseF = floorB
			}
			baseC := ceilA
			if !parity {
				baseC = ceilB
			}

			dist := math.Hypot(wx-g.p.pos.x, wy-g.p.pos.y)
			fog := clamp01(dist / maxDepth)
			ff := shade(baseF, 1.0-fog*0.7)
			cc := shade(baseC, 1.0-fog*0.7)

			if row > 0 {
				dst.Set(sx, sy, ff)
				if dist < g.zbuf[sx] {
					g.zbuf[sx] = dist
				}
			} else {
				dst.Set(sx, sy, cc)
			}

			floorX += stepX
			floorY += stepY
		}
	}
}
