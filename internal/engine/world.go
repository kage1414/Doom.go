package engine

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
)

func NewGame() *Game {
	w, h := DefaultMapW, DefaultMapH
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	grid, spawn, enemies, pickups := generateMap(w, h, rng)

	g := &Game{
		mapW:    w,
		mapH:    h,
		world:   grid,
		p:       player{pos: spawn, angle: 0, hp: playerStartHP, ammo: playerStartAmmo},
		state:   statePlaying,
		minimap: true,
		face:    basicfont.Face7x13,
	}
	g.fb = ebiten.NewImage(renderW, renderH)
	g.pix = ebiten.NewImage(1, 1)
	g.pix.Fill(white)
	g.scaleX = float64(ScreenW) / float64(renderW)
	g.scaleY = float64(ScreenH) / float64(renderH)
	g.zbuf = make([]float64, renderW)
	g.enemies = enemies
	g.pickups = pickups
	g.minimap = minimapOnAtStart
	g.mouseGrabbed = true

	sx, sy := int(math.Floor(spawn.x)), int(math.Floor(spawn.y))
	g.reachable = floodFillReachable(g.world, g.mapW, g.mapH, sx, sy)

	return g
}

func (g *Game) reset() {
	ng := NewGame()
	*g = *ng
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

// floodFillReachable marks all empty cells reachable from (sx,sy).
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
