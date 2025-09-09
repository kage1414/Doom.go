package engine

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type hitInfo struct {
	dist float64
	side int
}

func (g *Game) drawWalls(dst *ebiten.Image) {
	fov := deg2rad(fovDegrees)
	halfFov := fov / 2.0
	for x := 0; x < renderW; x++ {
		alpha := (float64(x)/float64(renderW))*fov - halfFov
		ray := normalizeAngle(g.p.angle + alpha)
		h := g.castRay(ray)

		corrected := h.dist * math.Cos(alpha)
		if corrected <= 0.0001 {
			corrected = 0.0001
		}
		g.zbuf[x] = corrected

		lineH := int(float64(renderH) / corrected)
		if lineH > renderH {
			lineH = renderH
		}
		start := renderH/2 - lineH/2

		var base color.RGBA
		switch h.side {
		case 0:
			base = colEast
		case 1:
			base = colWest
		case 2:
			base = colNorth
		default:
			base = colSouth
		}
		fog := clamp01(corrected / maxDepth)
		col := shade(base, 1.0-fog*0.85)

		drawRect(dst, g.pix, x, start, 1, lineH, col)
	}
}

func (g *Game) castRay(angle float64) hitInfo {
	sinA := math.Sin(angle)
	cosA := math.Cos(angle)
	mapX := int(math.Floor(g.p.pos.x))
	mapY := int(math.Floor(g.p.pos.y))

	var stepX, stepY int
	var sideDistX, sideDistY float64
	deltaDistX := math.Abs(1 / cosA)
	deltaDistY := math.Abs(1 / sinA)
	if math.IsInf(deltaDistX, 0) {
		deltaDistX = 1e30
	}
	if math.IsInf(deltaDistY, 0) {
		deltaDistY = 1e30
	}
	if cosA < 0 {
		stepX = -1
		sideDistX = (g.p.pos.x - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX+1) - g.p.pos.x) * deltaDistX
	}
	if sinA < 0 {
		stepY = -1
		sideDistY = (g.p.pos.y - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY+1) - g.p.pos.y) * deltaDistY
	}

	h := hitInfo{dist: 0, side: -1}
	for i := 0; i < 4096; i++ {
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
			h.side = mapSideVertical(stepX)
		} else {
			sideDistY += deltaDistY
			mapY += stepY
			h.side = mapSideHorizontal(stepY)
		}
		if mapX < 0 || mapY < 0 || mapX >= g.mapW || mapY >= g.mapH {
			h.dist = maxDepth
			break
		}
		t := g.world[mapY*g.mapW+mapX]
		if t == tWall {
			if h.side == 0 || h.side == 1 {
				h.dist = (float64(mapX) - g.p.pos.x + (1.0 - float64((stepX+1)/2))) / cosA
			} else {
				h.dist = (float64(mapY) - g.p.pos.y + (1.0 - float64((stepY+1)/2))) / sinA
			}
			if h.dist < 0.0001 {
				h.dist = 0.0001
			}
			if h.dist > maxDepth {
				h.dist = maxDepth
			}
			break
		}
	}
	return h
}
