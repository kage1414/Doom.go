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
		state:       stateStart,
		face:        basicfont.Face7x13,
		level:       1,
		totalLevels: DefaultLevels,
		minimap:     true,
	}
	g.fb = ebiten.NewImage(renderW, renderH)
	g.pix = ebiten.NewImage(1, 1)
	g.pix.Fill(white)
	g.scaleX = float64(ScreenW) / float64(renderW)
	g.scaleY = float64(ScreenH) / float64(renderH)
	g.zbuf = make([]float64, renderW)
	g.mouseGrabbed = false // start screen: mouse free

	g.initTextures()
	return g
}

// Piecewise scale: level 1 => 0.5x, middle => 1.0x, last => 3.0x
func scaleForLevel(level, total int) float64 {
	if total <= 1 {
		return 1.0
	}
	mid := (total + 1) / 2 // middle index (1-based)
	if level <= mid {
		// 1..mid maps 0.5 -> 1.0
		t := float64(level-1) / float64(mid-1)
		if mid == 1 {
			t = 1
		}
		return 0.5 + t*(1.0-0.5)
	}
	// mid..total maps 1.0 -> 3.0
	t := float64(level-mid) / float64(total-mid)
	if total == mid {
		t = 1
	}
	return 1.0 + t*(3.0-1.0)
}

// Randomize within ±pct (e.g., pct=0.30 => ±30%)
func jitter(val float64, pct float64, rng *rand.Rand) float64 {
	delta := (rng.Float64()*2 - 1) * pct
	return val * (1 + delta)
}

// clamp ints to >=1
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// setupLevel now uses piecewise scaling + jitter for map dims, enemies, and food
func (g *Game) setupLevel(level int, fresh bool) {
	if level < 1 {
		level = 1
	}
	if level > g.totalLevels {
		level = g.totalLevels
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Map dimensions
	scale := scaleForLevel(level, g.totalLevels)
	targetW := jitter(float64(MaxMapW)*scale, 0.30, rng)
	targetH := jitter(float64(MaxMapH)*scale, 0.30, rng)
	w := maxInt(int(targetW+0.5), BaseMapW/2) // keep reasonable minimums
	h := maxInt(int(targetH+0.5), BaseMapH/2)

	// Enemy total (we'll split by type later)
	targetEnemies := jitter(float64(BaseEnemyValue)*scale, 0.30, rng)
	totalEnemies := maxInt(int(targetEnemies+0.5), 1)

	// Food total (medkits + ammo)
	targetFood := jitter(float64(BaseFoodValue)*scale, 0.30, rng)
	totalFood := maxInt(int(targetFood+0.5), 1)

	// Split enemies by type (Zombies 60%, Runners 25%, Shooters 15%)
	ez := int(float64(totalEnemies) * 0.60)
	er := int(float64(totalEnemies) * 0.25)
	es := totalEnemies - ez - er
	if ez < 0 {
		ez = 0
	}
	if er < 0 {
		er = 0
	}
	if es < 0 {
		es = 0
	}

	// Split food 50/50 into medkits and ammo
	med := totalFood / 2
	amm := totalFood - med

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
	if g.level >= g.totalLevels {
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
