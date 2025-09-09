package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type vec2 struct{ x, y float64 }

type player struct {
	pos        vec2
	angle      float64
	hp         int
	ammo       int
	cooldown   float64
	muzzleTime float64
	score      int
}

type enemyType int

const (
	eZombie enemyType = iota
	eRunner
	eShooter
)

type enemy struct {
	pos    vec2
	hp     int
	etype  enemyType
	dead   bool
	blink  float64
	aiTime float64
}

type pickupType int

const (
	pickupMedkit pickupType = iota
	pickupAmmo
)

type pickup struct {
	pos   vec2
	ptype pickupType
	took  bool
}

type projectile struct {
	pos      vec2
	vel      vec2
	ttl      float64
	friendly bool
	radius   float64
	damage   int
}

type gameState int

const (
	stateStart gameState = iota
	statePlaying
	stateMenu
	stateLevelClear
	stateGameOver
	stateWin
)

const (
	tEmpty = 0
	tWall  = 1
)

type Game struct {
	mapW, mapH int
	world      []int
	reachable  []bool

	p       player
	enemies []*enemy
	pickups []*pickup
	bullets []*projectile

	level           int
	totalLevels     int
	defeated        int
	levelEnemyTotal int

	fb     *ebiten.Image
	pix    *ebiten.Image
	scaleX float64
	scaleY float64
	zbuf   []float64

	wallTex *ebiten.Image
	texW    int
	texH    int

	state        gameState
	minimap      bool
	mouseGrabbed bool
	lastMouseX   int

	face font.Face
}

var _ ebiten.Game = (*Game)(nil)
