package engine

import "image/color"

const (
	ScreenW = 960
	ScreenH = 600

	renderW = 320
	renderH = 200

	WallScale = 0.65

	// Level system
	LevelMax = 5
	BaseMapW = 48
	BaseMapH = 36
	MaxMapW  = 96
	MaxMapH  = 72

	// Generation controls
	MaxRooms        = 32
	RoomMinSize     = 4
	RoomMaxSize     = 10
	MinRoomSpacing  = 1
	SpawnSafeRadius = 6.0

	fovDegrees = 75.0
	moveSpeed  = 3.2
	sprintMul  = 1.8
	rotSpeed   = 2.6
	mouseSens  = 0.002

	maxDepth = 32.0

	minimapOnAtStart = true

	shootCooldownSec = 0.08
	playerMaxHP      = 100
	playerStartHP    = 85
	playerStartAmmo  = 36
	medkitHeal       = 25
	ammoPickupAmt    = 32

	playerShotSpd = 22.0
	playerShotTTL = 1.0
	playerShotDmg = 1

	zombieSpeed  = 1.35
	zombieHP     = 3
	runnerSpeed  = 2.25
	runnerHP     = 2
	shooterSpeed = 1.15
	shooterHP    = 3

	zombieRadius = 0.25
	touchDPS     = 10.0

	enemyShotCD  = 1.6
	enemyShotSpd = 12.0
	enemyShotDmg = 12
	enemyShotTTL = 1.6

	enemyKeepNear = 4.5
)

var (
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
