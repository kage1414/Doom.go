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

	for i, e := range g.enemies {
		if e.dead {
			continue
		}
		dx := e.pos.x - g.p.pos.x
		dy := e.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spriteEnemy, idx: i, dist: math.Hypot(dx, dy)})
	}

	for i, pk := range g.pickups {
		if pk.took {
			continue
		}
		dx := pk.pos.x - g.p.pos.x
		dy := pk.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spritePickup, idx: i, dist: math.Hypot(dx, dy)})
	}

	for i := range g.bullets {
		b := g.bullets[i]
		dx := b.pos.x - g.p.pos.x
		dy := b.pos.y - g.p.pos.y
		refs = append(refs, spriteRef{kind: spriteBullet, idx: i, dist: math.Hypot(dx, dy)})
	}

	sort.Slice(refs, func(i, j int) bool { return refs[i].dist > refs[j].dist })

	fov := deg2rad(fovDegrees)
	centerY := renderH / 2

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

			yTop := centerY - size/2
			headH := int(float64(size) * 0.3)
			bodyH := size - headH

			for x := startX; x <= endX; x++ {
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
				if r.dist <= g.zbuf[x] {
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
			y := centerY - size/2

			// Draw different sprites based on pickup type
			if pk.ptype == pickupAmmo {
				// Draw bullet-like shape
				g.drawBulletSprite(dst, startX, endX, y, size, dist)
			} else {
				// Draw first aid kit
				g.drawMedkitSprite(dst, startX, endX, y, size, dist)
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
			y := centerY - size/2
			for x := startX; x <= endX; x++ {
				if dist > g.zbuf[x] {
					continue
				}
				drawRect(dst, g.pix, x, y, 1, size, c)
			}
		}
	}
}

// drawBulletSprite draws a realistic bullet pickup sprite with animation
func (g *Game) drawBulletSprite(dst *ebiten.Image, startX, endX, y, size int, dist float64) {
	// Bullet colors - realistic bullet appearance
	bulletTip := color.RGBA{255, 255, 255, 255}  // White cone tip
	bulletBody := color.RGBA{200, 200, 200, 255} // Light grey cylinder
	bulletBase := color.RGBA{100, 100, 100, 255} // Dark grey base
	bulletRim := color.RGBA{80, 80, 80, 255}     // Dark rim

	// Animation: bobbing up and down
	bobOffset := int(math.Sin(g.gameTime*3.0) * 3.0) // 3 pixel bob, 3 cycles per second
	animatedY := y + bobOffset

	// Make the bullet skinny - reduce width to 1/3 of original
	bulletWidth := (endX - startX) / 3
	if bulletWidth < 1 {
		bulletWidth = 1
	}
	bulletStartX := startX + (endX-startX-bulletWidth)/2
	bulletEndX := bulletStartX + bulletWidth

	// Spinning animation: rotate the bullet around its center
	spinAngle := g.gameTime * 2.0 // 2 radians per second

	// Draw the bullet vertically - each horizontal slice represents a vertical section
	for x := bulletStartX; x <= bulletEndX; x++ {
		// Ensure x is within zbuf bounds
		if x < 0 || x >= len(g.zbuf) {
			continue
		}
		if dist > g.zbuf[x] {
			continue
		}

		// Calculate spinning offset for subtle animation
		spinOffset := math.Sin(spinAngle+float64(x-bulletStartX)*math.Pi*2/float64(bulletEndX-bulletStartX)) * 0.2

		// Draw the bullet vertically from top to bottom
		// Top 40% - white cone tip (pointing upward)
		coneHeight := int(float64(size) * 0.4)
		for y := 0; y < coneHeight; y++ {
			coneProgress := float64(y) / float64(coneHeight) // 0 at tip, 1 at base
			coneWidth := coneProgress                        // 0.0 at tip, 1.0 at base
			if float64(x-bulletStartX)/float64(bulletEndX-bulletStartX) <= coneWidth {
				bulletY := animatedY + y + int(spinOffset)
				drawRect(dst, g.pix, x, bulletY, 1, 1, bulletTip)
			}
		}

		// Next 50% - light grey body (cylinder)
		bodyHeight := int(float64(size) * 0.5)
		bodyStartY := coneHeight
		for y := 0; y < bodyHeight; y++ {
			bulletY := animatedY + bodyStartY + y + int(spinOffset)
			drawRect(dst, g.pix, x, bulletY, 1, 1, bulletBody)
		}

		// Bottom 10% - dark grey base (cylinder)
		baseHeight := int(float64(size) * 0.1)
		baseStartY := coneHeight + bodyHeight
		for y := 0; y < baseHeight; y++ {
			bulletY := animatedY + baseStartY + y + int(spinOffset)
			drawRect(dst, g.pix, x, bulletY, 1, 1, bulletBase)
		}

		// Add metallic rim effect
		if x == bulletStartX || x == bulletEndX {
			drawRect(dst, g.pix, x, animatedY, 1, size, bulletRim)
		}
	}
}

// drawMedkitSprite draws a realistic first aid kit pickup sprite with animation
func (g *Game) drawMedkitSprite(dst *ebiten.Image, startX, endX, y, size int, dist float64) {
	// First aid kit colors
	kitBody := color.RGBA{255, 255, 255, 255}   // White
	kitCross := color.RGBA{200, 50, 50, 255}    // Red cross
	kitBorder := color.RGBA{180, 180, 180, 255} // Gray border
	kitShadow := color.RGBA{220, 220, 220, 255} // Light shadow

	// Animation: bobbing up and down (same as bullets)
	bobOffset := int(math.Sin(g.gameTime*3.0) * 3.0) // 3 pixel bob, 3 cycles per second
	animatedY := y + bobOffset

	// Calculate medkit dimensions based on 3x3x1 ratio
	// The sprite should be roughly square (3x3) with some depth (1)
	medkitHeight := size

	// Draw the medkit as a 3D box with proper proportions
	for x := startX; x <= endX; x++ {
		// Ensure x is within zbuf bounds
		if x < 0 || x >= len(g.zbuf) {
			continue
		}
		if dist > g.zbuf[x] {
			continue
		}

		// Calculate relative position within sprite (0.0 to 1.0)
		relPos := float64(x-startX) / float64(endX-startX)

		// Create 3D box effect with depth
		var spriteHeight int
		var colorBody color.Color

		// Front face of the box (main body)
		if relPos < 0.7 {
			// Main body - full height
			spriteHeight = medkitHeight
			colorBody = kitBody
		} else {
			// Side face - slightly shorter for 3D effect
			spriteHeight = int(float64(medkitHeight) * 0.8)
			colorBody = kitShadow
		}

		// Draw the medkit body
		kitY := animatedY + (medkitHeight-spriteHeight)/2
		drawRect(dst, g.pix, x, kitY, 1, spriteHeight, colorBody)

		// Add red cross in the center (only on front face)
		if relPos >= 0.2 && relPos <= 0.5 {
			// Vertical cross line
			crossY := animatedY + medkitHeight/4
			crossH := medkitHeight / 2
			drawRect(dst, g.pix, x, crossY, 1, crossH, kitCross)

			// Horizontal cross line
			crossX := x
			crossW := 1
			if x > startX && x < endX {
				crossW = 3
				crossX = x - 1
			}
			centerY := animatedY + medkitHeight/2
			drawRect(dst, g.pix, crossX, centerY-1, crossW, 3, kitCross)
		}

		// Add border for definition
		if x == startX || x == endX {
			drawRect(dst, g.pix, x, animatedY, 1, medkitHeight, kitBorder)
		}
	}
}
