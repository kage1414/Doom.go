package engine

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type hitInfo struct {
	dist float64
	side int
	hx   float64
	hy   float64
}

func (g *Game) drawWalls(dst *ebiten.Image) {
	fov := deg2rad(fovDegrees)
	halfFov := fov / 2.0

	for x := 0; x < renderW; x++ {
		alpha := (float64(x)/float64(renderW))*fov - halfFov
		rayAng := normalizeAngle(g.p.angle + alpha)
		h := g.castRay(rayAng)

		corrected := h.dist * math.Cos(alpha)
		if corrected <= 0.0001 {
			corrected = 0.0001
		}
		g.zbuf[x] = corrected

		lineH := int(float64(renderH) / corrected * WallScale)
		start := renderH/2 - lineH/2

		var txf float64
		sinA := math.Sin(rayAng)
		cosA := math.Cos(rayAng)
		if h.side == 0 || h.side == 1 {
			txf = h.hy - math.Floor(h.hy)
			if cosA > 0 {
				txf = 1 - txf
			}
		} else {
			txf = h.hx - math.Floor(h.hx)
			if sinA < 0 {
				txf = 1 - txf
			}
		}
		tx := int(txf * float64(g.texW))
		if tx < 0 {
			tx = 0
		}
		if tx >= g.texW {
			tx = g.texW - 1
		}

		src := g.wallTex.SubImage(image.Rect(tx, 0, tx+1, g.texH)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1, float64(lineH)/float64(g.texH))
		op.GeoM.Translate(float64(x), float64(start))

		sideShade := 1.0
		if h.side == 0 || h.side == 2 {
			sideShade = 0.85
		}
		fog := clamp01(corrected / maxDepth)
		brightness := clamp01(sideShade * (1.0 - fog*0.85))
		op.ColorM.Scale(brightness, brightness, brightness, 1)
		dst.DrawImage(src, op)
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
		if g.world[mapY*g.mapW+mapX] == tWall {
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
			h.hx = g.p.pos.x + cosA*h.dist
			h.hy = g.p.pos.y + sinA*h.dist
			break
		}
	}
	return h
}
