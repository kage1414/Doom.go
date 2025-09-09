package engine

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Layout(outW, outH int) (int, int) {
	if outW > 0 && outH > 0 {
		g.scaleX = float64(outW) / float64(renderW)
		g.scaleY = float64(outH) / float64(renderH)
	}
	return ScreenW, ScreenH
}

func (g *Game) Update() error {
	// Global Esc behavior
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		switch g.state {
		case statePlaying:
			g.state = stateMenu
			g.mouseGrabbed = false
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		case stateMenu:
			return ebiten.Termination
		}
	}

	switch g.state {
	case stateStart:
		// Choose total levels before starting
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if g.totalLevels < MaxLevelCap {
				g.totalLevels++
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if g.totalLevels > 1 {
				g.totalLevels--
			}
		}
		// Digit keys quick-set
		for key := ebiten.Key0; key <= ebiten.Key9; key++ {
			if inpututil.IsKeyJustPressed(key) {
				n := int(key - ebiten.Key0)
				if n == 0 {
					n = 10
				}
				if n > MaxLevelCap {
					n = MaxLevelCap
				}
				g.totalLevels = n
			}
		}
		// Enter to begin
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.level = 1
			g.setupLevel(g.level, true)
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		return nil

	case stateMenu:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			return ebiten.Termination
		}
		return nil

	case stateLevelClear:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.level++
			if g.level > g.totalLevels {
				g.state = stateWin
				return nil
			}
			g.setupLevel(g.level, false)
			g.state = statePlaying
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			g.lastMouseX = 0
		}
		return nil

	case stateGameOver, stateWin:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
		return nil

	case statePlaying:
		dt := 1.0 / 60.0

		if g.p.cooldown > 0 {
			g.p.cooldown -= dt
			if g.p.cooldown < 0 {
				g.p.cooldown = 0
			}
		}
		if g.p.muzzleTime > 0 {
			g.p.muzzleTime -= dt
			if g.p.muzzleTime < 0 {
				g.p.muzzleTime = 0
			}
		}

		for _, e := range g.enemies {
			if e.blink > 0 {
				e.blink -= dt
				if e.blink < 0 {
					e.blink = 0
				}
			}
			if !e.dead {
				e.aiTime += dt
			}
		}

		g.updateProjectiles(dt)

		if inpututil.IsKeyJustPressed(ebiten.KeyM) {
			g.minimap = !g.minimap
		}

		if !g.mouseGrabbed {
			g.mouseGrabbed = true
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
		}

		if g.mouseGrabbed {
			x, _ := ebiten.CursorPosition()
			if g.lastMouseX != 0 {
				dx := x - g.lastMouseX
				g.p.angle += float64(dx) * mouseSens
				g.p.angle = normalizeAngle(g.p.angle)
			}
			g.lastMouseX = x
		} else {
			g.lastMouseX = 0
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.p.angle -= rotSpeed * dt
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.p.angle += rotSpeed * dt
		}
		g.p.angle = normalizeAngle(g.p.angle)

		forward, side := 0.0, 0.0
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			forward++
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) {
			forward--
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			side--
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			side++
		}
		speed := moveSpeed
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			speed *= sprintMul
		}
		if forward != 0 || side != 0 {
			l := math.Hypot(forward, side)
			forward /= l
			side /= l
			fx := math.Cos(g.p.angle)
			fy := math.Sin(g.p.angle)
			rx := -fy
			ry := fx
			vx := (fx*forward + rx*side) * speed * dt
			vy := (fy*forward + ry*side) * speed * dt
			g.moveWithCollision(vx, vy)
		}

		if (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace)) &&
			g.p.cooldown <= 0 && g.p.ammo > 0 {
			g.p.cooldown = shootCooldownSec
			g.p.muzzleTime = 0.06
			g.p.ammo--
			g.firePlayerShot()
		}

		for _, e := range g.enemies {
			if e.dead {
				continue
			}
			switch e.etype {
			case eZombie:
				g.seekEnemy(e, zombieSpeed, dt)
				if dist2(e.pos.x, e.pos.y, g.p.pos.x, g.p.pos.y) < (0.25+0.25)*(0.25+0.25) {
					g.p.hp -= int(touchDPS * dt)
					if g.p.hp < 0 {
						g.p.hp = 0
					}
				}
			case eRunner:
				g.seekEnemy(e, runnerSpeed, dt)
				if dist2(e.pos.x, e.pos.y, g.p.pos.x, g.p.pos.y) < (0.25+0.25)*(0.25+0.25) {
					g.p.hp -= int(touchDPS * dt)
					if g.p.hp < 0 {
						g.p.hp = 0
					}
				}
			case eShooter:
				g.shooterAI(e, dt)
			}
		}

		for _, pk := range g.pickups {
			if pk.took {
				continue
			}
			if dist2(pk.pos.x, pk.pos.y, g.p.pos.x, g.p.pos.y) < 0.5*0.5 {
				switch pk.ptype {
				case pickupMedkit:
					if g.p.hp < playerMaxHP {
						g.p.hp += medkitHeal
						if g.p.hp > playerMaxHP {
							g.p.hp = playerMaxHP
						}
						pk.took = true
						// Add pickup message
						g.pickupMessages = append(g.pickupMessages, pickupMessage{
							text:     fmt.Sprintf("+%d Health", medkitHeal),
							color:    green,
							timeLeft: pickupMessageDuration,
						})
					}
				case pickupAmmo:
					g.p.ammo += ammoPickupAmt
					pk.took = true
					// Add pickup message
					g.pickupMessages = append(g.pickupMessages, pickupMessage{
						text:     fmt.Sprintf("+%d Ammo", ammoPickupAmt),
						color:    yellow,
						timeLeft: pickupMessageDuration,
					})
				}
			}
		}

		if g.p.hp <= 0 {
			g.state = stateGameOver
			g.mouseGrabbed = false
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
			return nil
		}

		allDead := true
		for _, e := range g.enemies {
			if !e.dead {
				allDead = false
				break
			}
		}
		if allDead {
			g.advanceLevelOrWin()
			return nil
		}
		
		// Update pickup messages
		g.updatePickupMessages(dt)
	}
	return nil
}

func (g *Game) moveWithCollision(dx, dy float64) {
	newX := g.p.pos.x + dx
	newY := g.p.pos.y + dy
	if !g.isSolid(int(math.Floor(newX)), int(math.Floor(g.p.pos.y))) {
		g.p.pos.x = newX
	}
	if !g.isSolid(int(math.Floor(g.p.pos.x)), int(math.Floor(newY))) {
		g.p.pos.y = newY
	}
}

func (g *Game) updatePickupMessages(dt float64) {
	// Update message timers and remove expired messages
	nm := g.pickupMessages[:0]
	for _, msg := range g.pickupMessages {
		msg.timeLeft -= dt
		if msg.timeLeft > 0 {
			nm = append(nm, msg)
		}
	}
	g.pickupMessages = nm
}
