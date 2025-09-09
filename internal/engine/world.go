package engine

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
)

func NewGame() *Game {
	g := &Game{
		state:   statePlaying,
		face:    basicfont.Face7x13,
		level:   1,
		minimap: true,
	}
	g.fb = ebiten.NewImage(renderW, renderH)
	g.pix = ebiten.NewImage(1, 1)
	g.pix.Fill(white)
	g.scaleX = float64(ScreenW) / float64(renderW)
	g.scaleY = float64(ScreenH) / float64(renderH)
	g.zbuf = make([]float64, renderW)
	g.mouseGrabbed = true

	g.initTextures()
	g.setupLevel(g.level, true)
	return g
}

func (g *Game) setupLevel(level int, fresh bool) {
	lerp := func(a, b, t float64) int { return int(a + (b-a)*t + 0.5) }
	if level < 1 {
		level = 1
	}
	if level > LevelMax {
		level = LevelMax
	}
	t := float64(level-1) / float64(LevelMax-1)

	w := lerp(float64(BaseMapW), float64(MaxMapW), t)
	h := lerp(float64(BaseMapH), float64(MaxMapH), t)

	baseZ, maxZ := 5, 20
	baseR, maxR := 2, 12
	baseS, maxS := 1, 10
	ez := lerp(float64(baseZ), float64(maxZ), t)
	er := lerp(float64(baseR), float64(maxR), t)
	es := lerp(float64(baseS), float64(maxS), t)

	baseMed, minMed := 10, 4
	baseAmmo, minAmmo := 12, 6
	med := lerp(float64(baseMed), float64(minMed), t)
	amm := lerp(float64(baseAmmo), float64(minAmmo), t)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	grid, spawn, enemies, pickups := generateMap(w, h, rng, ez, er, es, med, amm)

	g.mapW, g.mapH = w, h
	g.world = grid
	g.enemies = enemies
	g.pickups = pickups
	g.levelEnemyTotal = len(enemies)

	if fresh {
		g.p = player{pos: spawn, angle: 0, hp: playerStartHP, ammo: playerStartAmmo}
		g.defeated = 0
	} else {
		g.p.pos = spawn
		g.p.angle = 0
	}

	sx, sy := int(math.Floor(spawn.x)), int(math.Floor(spawn.y))
	g.reachable = floodFillReachable(g.world, g.mapW, g.mapH, sx, sy)
}

func (g *Game) reset() {
	ng := NewGame()
	*g = *ng
}

func (g *Game) advanceLevelOrWin() {
	if g.level >= LevelMax {
		g.state = stateWin
		g.mouseGrabbed = false
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
		return
	}
	g.state = stateLevelClear
	g.mouseGrabbed = false
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

func (g *Game) isSolid(ix, iy int) bool {
	if ix < 0 || iy < 0 || ix >= g.mapW || iy >= g.mapH {
		return true
	}
	return g.world[iy*g.mapW+ix] == tWall
}

func (g *Game) isSolidAtFloat(x, y float64) bool {
	return g.isSolid(int(math.Floor(x)), int(math.Floor(y)))
}

func floodFillReachable(grid []int, w, h, sx, sy int) []bool {
	reach := make([]bool, w*h)
	if sx < 0 || sy < 0 || sx >= w || sy >= h {
		return reach
	}
	if grid[sy*w+sx] == tWall {
		return reach
	}
	qx := make([]int, 0, w*h/4)
	qy := make([]int, 0, w*h/4)
	push := func(x, y int) {
		idx := y*w + x
		if x < 0 || y < 0 || x >= w || y >= h {
			return
		}
		if grid[idx] == tWall || reach[idx] {
			return
		}
		reach[idx] = true
		qx = append(qx, x)
		qy = append(qy, y)
	}
	push(sx, sy)
	head := 0
	for head < len(qx) {
		cx, cy := qx[head], qy[head]
		head++
		push(cx+1, cy)
		push(cx-1, cy)
		push(cx, cy+1)
		push(cx, cy-1)
	}
	return reach
}
