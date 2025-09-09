package engine

import "image/color"

const (
	ScreenW = 960
	ScreenH = 600

	renderW = 320
	renderH = 200

	// Procedural map size
	DefaultMapW = 64
	DefaultMapH = 48

	// Generation controls
	MaxRooms        = 22
	RoomMinSize     = 4
	RoomMaxSize     = 10
	MinRoomSpacing  = 1   // gap between rooms
	SpawnSafeRadius = 6.0 // tiles around spawn without enemies

	EnemiesZombie  = 10
	EnemiesRunner  = 6
	EnemiesShooter = 6

	PickupsMedkit = 8
	PickupsAmmo   = 10

	fovDegrees = 75.0
	moveSpeed  = 3.2
	sprintMul  = 1.8
	rotSpeed   = 2.6
	mouseSens  = 0.002
	maxDepth   = 32.0

	minimapOnAtStart = true

	// fire rate
	shootCooldownSec = 0.08
	playerMaxHP      = 100
	playerStartHP    = 85
	playerStartAmmo  = 360
	medkitHeal       = 25
	ammoPickupAmt    = 32
	playerShotSpd    = 10.0
	playerShotTTL    = 1.0
	playerShotDmg    = 1

	zombieSpeed  = 1.35
	zombieHP     = 3
	runnerSpeed  = 2.25
	runnerHP     = 2
	shooterSpeed = 1.15
	shooterHP    = 3

	zombieRadius = 0.25
	touchDPS     = 10.0

	enemyShotCD  = 1.6
	enemyShotSpd = 6.0
	enemyShotDmg = 12
	enemyShotTTL = 1.6

	enemyKeepNear = 4.5
)

var (
	colNorth = color.RGBA{180, 30, 30, 255}
	colSouth = color.RGBA{220, 50, 50, 255}
	colWest  = color.RGBA{30, 30, 180, 255}
	colEast  = color.RGBA{50, 50, 220, 255}

	floorA = color.RGBA{26, 28, 26, 255}
	floorB = color.RGBA{36, 40, 36, 255}
	ceilA  = color.RGBA{10, 12, 16, 255}
	ceilB  = color.RGBA{14, 16, 20, 255}

	uiBox    = color.RGBA{0, 0, 0, 200}
	uiAccent = color.RGBA{200, 255, 200, 255}
	white    = color.RGBA{255, 255, 255, 255}
	yellow   = color.RGBA{240, 220, 120, 255}
	green    = color.RGBA{110, 200, 120, 255}
	red      = color.RGBA{230, 60, 60, 255}
	gray     = color.RGBA{150, 150, 150, 255}
	cyan     = color.RGBA{120, 210, 230, 255}
	magenta  = color.RGBA{210, 120, 230, 255}
	black    = color.RGBA{0, 0, 0, 255}
)
