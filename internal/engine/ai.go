package engine

import "math"

const (
	stepSize    = 0.08
	enemyRadius = 0.25 // radius for collision checks
)

// Move an enemy with circle-vs-grid collision using swept steps.
func (g *Game) moveEnemyCircle(e *enemy, dx, dy, radius float64) {
	steps := int(math.Ceil(math.Max(math.Abs(dx), math.Abs(dy)) / stepSize))
	if steps < 1 {
		steps = 1
	}
	sx := dx / float64(steps)
	sy := dy / float64(steps)
	for i := 0; i < steps; i++ {
		nx := e.pos.x + sx
		if !g.circleHitsSolid(nx, e.pos.y, radius) {
			e.pos.x = nx
		}
		ny := e.pos.y + sy
		if !g.circleHitsSolid(e.pos.x, ny, radius) {
			e.pos.y = ny
		}
	}
}

// Circle vs grid check by sampling around the circle.
func (g *Game) circleHitsSolid(cx, cy, r float64) bool {
	if g.isSolidAtFloat(cx-r, cy) {
		return true
	}
	if g.isSolidAtFloat(cx+r, cy) {
		return true
	}
	if g.isSolidAtFloat(cx, cy-r) {
		return true
	}
	if g.isSolidAtFloat(cx, cy+r) {
		return true
	}
	return g.isSolidAtFloat(cx, cy)
}

func (g *Game) seekEnemy(e *enemy, speed, dt float64) {
	dx := g.p.pos.x - e.pos.x
	dy := g.p.pos.y - e.pos.y
	dist := math.Hypot(dx, dy)
	if dist < 1e-6 {
		return
	}
	vx := (dx / dist) * speed * dt
	vy := (dy / dist) * speed * dt
	g.moveEnemyCircle(e, vx, vy, enemyRadius)
}

func (g *Game) shooterAI(e *enemy, dt float64) {
	dx := g.p.pos.x - e.pos.x
	dy := g.p.pos.y - e.pos.y
	dist := math.Hypot(dx, dy)
	dirx := dx / (dist + 1e-6)
	diry := dy / (dist + 1e-6)

	desired := enemyKeepNear
	err := dist - desired
	move := clampF(err*0.7, -1.5, 1.5)
	strafe := 0.8
	tx := dirx*move - diry*strafe
	ty := diry*move + dirx*strafe

	g.moveEnemyCircle(e, tx*shooterSpeed*dt, ty*shooterSpeed*dt, enemyRadius)

	if e.aiTime >= enemyShotCD && g.hasLineOfSightGrid(e.pos, g.p.pos) {
		e.aiTime = 0
		v := vec2{dirx * enemyShotSpd, diry * enemyShotSpd}
		g.bullets = append(g.bullets, &projectile{
			pos:      vec2{e.pos.x + dirx*0.3, e.pos.y + diry*0.3},
			vel:      v,
			ttl:      enemyShotTTL,
			friendly: false,
			radius:   0.05,
			damage:   enemyShotDmg,
		})
	}
}

// DDA line-of-sight: only true if we reach the target cell before hitting a wall/closed door.
func (g *Game) hasLineOfSightGrid(a, b vec2) bool {
	ax, ay := a.x, a.y
	bx, by := b.x, b.y

	rayDirX := bx - ax
	rayDirY := by - ay
	dist := math.Hypot(rayDirX, rayDirY)
	if dist < 1e-6 {
		return true
	}
	rayDirX /= dist
	rayDirY /= dist

	mapX := int(math.Floor(ax))
	mapY := int(math.Floor(ay))
	endX := int(math.Floor(bx))
	endY := int(math.Floor(by))

	deltaDistX := math.Abs(1.0 / (rayDirX + 1e-12))
	deltaDistY := math.Abs(1.0 / (rayDirY + 1e-12))

	var stepX, stepY int
	var sideDistX, sideDistY float64

	if rayDirX < 0 {
		stepX = -1
		sideDistX = (ax - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX+1) - ax) * deltaDistX
	}
	if rayDirY < 0 {
		stepY = -1
		sideDistY = (ay - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY+1) - ay) * deltaDistY
	}

	for i := 0; i < 4096; i++ {
		if mapX == endX && mapY == endY {
			return true
		}
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
		} else {
			sideDistY += deltaDistY
			mapY += stepY
		}
		if mapX < 0 || mapY < 0 || mapX >= g.mapW || mapY >= g.mapH {
			return false
		}
		if g.isSolid(mapX, mapY) {
			return false
		}
	}
	return false
}
