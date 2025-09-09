package engine

import (
	"image/color"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

func enemyMaxHP(e *enemy) int {
	switch e.etype {
	case eZombie:
		return zombieHP
	case eRunner:
		return runnerHP
	case eShooter:
		return shooterHP
	default:
		return 1
	}
}

type spriteKind int

const (
	spriteEnemy spriteKind = iota
	spritePickup
	spriteBullet
)

type spriteRef struct {
	kind spriteKind
	idx  int
	dist float64
}

func (g *Game) drawSprites(dst *ebiten.Image) {
	refs := make([]spriteRef, 0, len(g.enemies)+len(g.pickups)+len(g.bullets))

	// collect enemies
	for i, e := range g.enemies {
		if e.dead {
			continue
		}
		dx := e.pos.x - g.p.pos.x
		dy := e.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spriteEnemy, idx: i, dist: math.Hypot(dx, dy)})
	}

	// collect pickups
	for i, pk := range g.pickups {
		if pk.took {
			continue
		}
		dx := pk.pos.x - g.p.pos.x
		dy := pk.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spritePickup, idx: i, dist: math.Hypot(dx, dy)})
	}

	// collect bullets
	for i := range g.bullets {
		b := g.bullets[i]
		dx := b.pos.x - g.p.pos.x
		dy := b.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spriteBullet, idx: i, dist: math.Hypot(dx, dy)})
	}

	// sort far -> near so nearer sprites overwrite farther ones
	sort.Slice(refs, func(i, j int) bool { return refs[i].dist > refs[j].dist })

	fov := deg2rad(fovDegrees)

	for _, r := range refs {
		switch r.kind {
		case spriteEnemy:
			e := g.enemies[r.idx]
			dx := e.pos.x - g.p.pos.x
			dy := e.pos.y - g.p.pos.y
			dist := math.Hypot(dx, dy)
			if dist <= 0.001 {
				continue
			}
			ang := math.Atan2(dy, dx) - g.p.angle
			ang = normalizeAngle(ang)
			if ang > math.Pi {
				ang -= 2 * math.Pi
			}
			if math.Abs(ang) > fov {
				continue
			}

			size := int(float64(renderH) / dist * 0.55)
			if size < 2 {
				size = 2
			}
			screenX := int((0.5 + (ang / fov)) * float64(renderW))
			startX := screenX - size/3
			endX := screenX + size/3
			if startX < 0 {
				startX = 0
			}
			if endX > renderW-1 {
				endX = renderW - 1
			}

			bodyCol := gray
			headCol := white
			switch e.etype {
			case eZombie:
				bodyCol = gray
				headCol = color.RGBA{210, 210, 210, 255}
			case eRunner:
				bodyCol = cyan
				headCol = color.RGBA{220, 240, 255, 255}
			case eShooter:
				bodyCol = magenta
				headCol = color.RGBA{250, 210, 255, 255}
			}
			if e.blink > 0 {
				bodyCol = white
				headCol = white
			}

			yTop := renderH/2 - size/2
			headH := int(float64(size) * 0.3)
			bodyH := size - headH

			for x := startX; x <= endX; x++ {
				// don't draw behind a nearer wall
				if dist > g.zbuf[x] {
					continue
				}
				y := yTop
				if headH > 0 {
					drawRect(dst, g.pix, x, y, 1, headH, headCol)
					y += headH
				}
				if bodyH > 0 {
					drawRect(dst, g.pix, x, y, 1, bodyH, bodyCol)
				}
				if x == startX || x == endX {
					drawRect(dst, g.pix, x, yTop, 1, size, color.RGBA{0, 0, 0, 120})
				}
			}

			// HP bar
			hpMax := enemyMaxHP(e)
			if hpMax < 1 {
				hpMax = 1
			}
			barW := endX - startX + 1
			if barW < 6 {
				barW = 6
			}
			barH := 2
			barX := startX
			barY := yTop - 4
			visible := false
			for x := startX; x <= endX && !visible; x++ {
				if dist <= g.zbuf[x] {
					visible = true
				}
			}
			if visible {
				drawRect(dst, g.pix, barX, barY, barW, barH, black)
				fillW := int(float64(barW) * clamp01(float64(e.hp)/float64(hpMax)))
				if fillW > 0 {
					col := red
					if e.hp >= (hpMax+1)/2 {
						col = green
					} else if e.hp > 1 {
						col = yellow
					}
					drawRect(dst, g.pix, barX, barY, fillW, barH, col)
				}
			}

		case spritePickup:
			pk := g.pickups[r.idx]
			dx := pk.pos.x - g.p.pos.x
			dy := pk.pos.y - g.p.pos.y
			dist := math.Hypot(dx, dy)
			if dist <= 0.001 {
				continue
			}
			ang := math.Atan2(dy, dx) - g.p.angle
			ang = normalizeAngle(ang)
			if ang > math.Pi {
				ang -= 2 * math.Pi
			}
			if math.Abs(ang) > fov {
				continue
			}
			size := int(float64(renderH) / dist * 0.35)
			if size < 1 {
				size = 1
			}
			screenX := int((0.5 + (ang / fov)) * float64(renderW))
			startX := screenX - size/2
			endX := screenX + size/2
			if startX < 0 {
				startX = 0
			}
			if endX > renderW-1 {
				endX = renderW - 1
			}
			colorBody := green
			if pk.ptype == pickupAmmo {
				colorBody = yellow
			}
			y := renderH/2 - size/2
			for x := startX; x <= endX; x++ {
				if dist > g.zbuf[x] {
					continue
				}
				drawRect(dst, g.pix, x, y, 1, size, colorBody)
				if x == startX || x == endX {
					drawRect(dst, g.pix, x, y, 1, size, color.RGBA{0, 0, 0, 120})
				}
			}

		case spriteBullet:
			b := g.bullets[r.idx]
			dx := b.pos.x - g.p.pos.x
			dy := b.pos.y - g.p.pos.y
			dist := math.Hypot(dx, dy)
			if dist <= 0.001 {
				continue
			}
			ang := math.Atan2(dy, dx) - g.p.angle
			ang = normalizeAngle(ang)
			if ang > math.Pi {
				ang -= 2 * math.Pi
			}
			if math.Abs(ang) > fov {
				continue
			}
			size := int(float64(renderH) / dist * 0.2)
			if size < 1 {
				size = 1
			}
			screenX := int((0.5 + (ang / fov)) * float64(renderW))
			startX := screenX - 1
			endX := screenX + 1
			if startX < 0 {
				startX = 0
			}
			if endX > renderW-1 {
				endX = renderW - 1
			}
			c := yellow
			if !b.friendly {
				c = red
			}
			y := renderH/2 - size/2
			for x := startX; x <= endX; x++ {
				// respect wall depth: don't draw behind nearer wall columns
				if dist > g.zbuf[x] {
					continue
				}
				drawRect(dst, g.pix, x, y, 1, size, c)
			}
		}
	}
}
