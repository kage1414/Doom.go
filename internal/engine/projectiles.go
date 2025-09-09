package engine

import "math"

func (g *Game) firePlayerShot() {
	dirx, diry := math.Cos(g.p.angle), math.Sin(g.p.angle)
	g.bullets = append(g.bullets, &projectile{
		pos:      vec2{g.p.pos.x + dirx*0.4, g.p.pos.y + diry*0.4},
		vel:      vec2{dirx * playerShotSpd, diry * playerShotSpd},
		ttl:      playerShotTTL,
		friendly: true,
		radius:   0.05,
		damage:   playerShotDmg,
	})
}

func (g *Game) updateProjectiles(dt float64) {
	nb := g.bullets[:0]
	for _, b := range g.bullets {
		b.ttl -= dt
		if b.ttl <= 0 {
			continue
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
						}
						b.ttl = 0
						goto bulletDone
					}
				}
			} else {
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
