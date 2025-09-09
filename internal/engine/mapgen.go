package engine

import (
	"math"
	"math/rand"
)

type rect struct{ x, y, w, h int }

func (r rect) center() (int, int) { return r.x + r.w/2, r.y + r.h/2 }
func (r rect) intersects(o rect, padding int) bool {
	return r.x < o.x+o.w+padding &&
		r.x+r.w+padding > o.x &&
		r.y < o.y+o.h+padding &&
		r.y+r.h+padding > o.y
}

// generateMap builds a room/corridor map and scatters enemies/pickups based on inputs.
func generateMap(w, h int, rng *rand.Rand, ez, er, es, medkits, ammos int) (grid []int, spawn vec2, enemies []*enemy, pickups []*pickup) {
	grid = make([]int, w*h)
	for i := range grid {
		grid[i] = tWall
	}

	rooms := make([]rect, 0, MaxRooms)

	tryRooms := MaxRooms * 3
	for i := 0; i < tryRooms && len(rooms) < MaxRooms; i++ {
		rw := rng.Intn(RoomMaxSize-RoomMinSize+1) + RoomMinSize
		rh := rng.Intn(RoomMaxSize-RoomMinSize+1) + RoomMinSize
		rx := rng.Intn(w-rw-2) + 1
		ry := rng.Intn(h-rh-2) + 1
		rc := rect{rx, ry, rw, rh}
		ok := true
		for _, other := range rooms {
			if rc.intersects(other, MinRoomSpacing) {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		digRoom(grid, w, h, rc)
		rooms = append(rooms, rc)

		if len(rooms) > 1 {
			x1, y1 := rooms[len(rooms)-2].center()
			x2, y2 := rc.center()
			if rng.Intn(2) == 0 {
				digH2(grid, w, h, x1, x2, y1)
				digV2(grid, w, h, y1, y2, x2)
			} else {
				digV2(grid, w, h, y1, y2, x1)
				digH2(grid, w, h, x1, x2, y2)
			}
		}
	}

	if len(rooms) == 0 {
		r := rect{w / 2, h / 2, 5, 5}
		digRoom(grid, w, h, r)
		digH2(grid, w, h, 2, w-3, r.y+r.h/2)
		digV2(grid, w, h, 2, h-3, r.x+r.w/2)
		rooms = append(rooms, r)
	}

	sx, sy := rooms[0].center()
	spawn = vec2{float64(sx) + 0.5, float64(sy) + 0.5}

	spreadEnemy := func(count int, kind enemyType, hp int) {
		for placed := 0; placed < count; {
			x := rng.Intn(w-2) + 1
			y := rng.Intn(h-2) + 1
			if grid[y*w+x] != tEmpty {
				continue
			}
			if math.Hypot(float64(x)-spawn.x+0.5, float64(y)-spawn.y+0.5) < SpawnSafeRadius {
				continue
			}
			enemies = append(enemies, &enemy{
				pos:   vec2{float64(x) + 0.5, float64(y) + 0.5},
				hp:    hp,
				etype: kind,
			})
			placed++
		}
	}
	spreadEnemy(ez, eZombie, zombieHP)
	spreadEnemy(er, eRunner, runnerHP)
	spreadEnemy(es, eShooter, shooterHP)

	placePickup := func(count int, pt pickupType) {
		for placed := 0; placed < count; {
			x := rng.Intn(w-2) + 1
			y := rng.Intn(h-2) + 1
			if grid[y*w+x] != tEmpty {
				continue
			}
			if math.Hypot(float64(x)-spawn.x+0.5, float64(y)-spawn.y+0.5) < 3.5 {
				continue
			}
			pickups = append(pickups, &pickup{
				pos:   vec2{float64(x) + 0.5, float64(y) + 0.5},
				ptype: pt,
			})
			placed++
		}
	}
	placePickup(medkits, pickupMedkit)
	placePickup(ammos, pickupAmmo)

	return grid, spawn, enemies, pickups
}

func digRoom(grid []int, w, h int, r rect) {
	for y := r.y; y < r.y+r.h; y++ {
		for x := r.x; x < r.x+r.w; x++ {
			if x > 0 && y > 0 && x < w-1 && y < h-1 {
				grid[y*w+x] = tEmpty
			}
		}
	}
}

// 2-wide corridors (already in place previously)
func digH2(grid []int, w, h, x1, x2, y int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for x := x1; x <= x2; x++ {
		if x <= 0 || x >= w-1 {
			continue
		}
		if y > 0 && y < h-1 {
			grid[y*w+x] = tEmpty
		}
		if y+1 > 0 && y+1 < h-1 {
			grid[(y+1)*w+x] = tEmpty
		} else if y-1 > 0 && y-1 < h-1 {
			grid[(y-1)*w+x] = tEmpty
		}
	}
}

func digV2(grid []int, w, h, y1, y2, x int) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		if y <= 0 || y >= h-1 {
			continue
		}
		if x > 0 && x < w-1 {
			grid[y*w+x] = tEmpty
		}
		if x+1 > 0 && x+1 < w-1 {
			grid[y*w+(x+1)] = tEmpty
		} else if x-1 > 0 && x-1 < w-1 {
			grid[y*w+(x-1)] = tEmpty
		}
	}
}
