package engine

import (
	"math"
	"math/rand"
)

func (g *Game) firePlayerShot() {
	dirx, diry := math.Cos(g.p.angle), math.Sin(g.p.angle)
	g.bullets = append(g.bullets, &projectile{
		pos:        vec2{g.p.pos.x + dirx*0.4, g.p.pos.y + diry*0.4},
		vel:        vec2{dirx * g.settings.bulletSpeed, diry * g.settings.bulletSpeed},
		ttl:        playerShotTTL,
		friendly:   true,
		radius:     0.05,
		damage:     playerShotDmg,
		curveAngle: (rand.Float64() - 0.5) * 0.3, // Random curve between -0.15 and 0.15 radians
		curveRate:  0.5 + rand.Float64()*0.5,     // Curve rate between 0.5 and 1.0
	})

	// Play bullet sound
	g.playBulletSound()
}

func (g *Game) updateProjectiles(dt float64) {
	nb := g.bullets[:0]
	for _, b := range g.bullets {
		b.ttl -= dt
		if b.ttl <= 0 {
			continue
		}

		// Apply curving effect - gradually change velocity direction
		// The curve becomes more pronounced as the bullet travels further
		curveStrength := (playerShotTTL - b.ttl) * b.curveRate * 0.05 // Increases over time

		// Apply curve perpendicular to current velocity direction
		velLen := math.Hypot(b.vel.x, b.vel.y)
		if velLen > 0 {
			// Get perpendicular direction (rotate 90 degrees)
			perpX := -b.vel.y / velLen
			perpY := b.vel.x / velLen

			// Apply curve in the perpendicular direction
			curveX := perpX * curveStrength * math.Cos(b.curveAngle)
			curveY := perpY * curveStrength * math.Sin(b.curveAngle)

			// Add curve to velocity
			b.vel.x += curveX * dt
			b.vel.y += curveY * dt
		}

		steps := int(math.Ceil(math.Max(math.Abs(b.vel.x*dt), math.Abs(b.vel.y*dt)) / 0.05))
		if steps < 1 {
			steps = 1
		}
		sx := (b.vel.x * dt) / float64(steps)
		sy := (b.vel.y * dt) / float64(steps)
		hitWall := false

		for i := 0; i < steps; i++ {
			nx := b.pos.x + sx
			ny := b.pos.y + sy
			if g.isSolidAtFloat(nx, ny) {
				hitWall = true
				break
			}
			b.pos.x, b.pos.y = nx, ny

			if b.friendly {
				for _, e := range g.enemies {
					if e.dead {
						continue
					}
					if dist2(b.pos.x, b.pos.y, e.pos.x, e.pos.y) < 0.35*0.35 {
						e.hp -= b.damage
						e.blink = 0.12
						if e.hp <= 0 {
							e.dead = true
							g.defeated++      // <- track defeated enemies
							g.playCoinSound() // Play coin sound when enemy dies
						}
						b.ttl = 0
						goto bulletDone
					}
				}
			} else {
				// Check for bullet whiz sound (enemy bullets passing close to player)
				if !b.whizPlayed {
					distToPlayer := math.Hypot(b.pos.x-g.p.pos.x, b.pos.y-g.p.pos.y)
					if distToPlayer < 1.5 && distToPlayer > 0.5 { // Close but not hitting
						g.playBulletWhizSound()
						b.whizPlayed = true
					}
				}

				if dist2(b.pos.x, b.pos.y, g.p.pos.x, g.p.pos.y) < 0.35*0.35 {
					g.p.hp -= b.damage
					if g.p.hp < 0 {
						g.p.hp = 0
					}
					b.ttl = 0
					goto bulletDone
				}
			}
		}

	bulletDone:
		if !hitWall && b.ttl > 0 {
			nb = append(nb, b)
		}
	}
	g.bullets = nb
}
