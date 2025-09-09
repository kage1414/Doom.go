package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"doomlike/internal/engine"
)

func main() {
	g := engine.NewGame()
	ebiten.SetWindowSize(engine.ScreenW, engine.ScreenH)
	ebiten.SetWindowTitle("DOOM.go — Sprites, Health Bars, Win/Lose (Esc: Menu)")
	ebiten.SetWindowResizable(true)
	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	if err := ebiten.RunGame(g); err != nil && err != ebiten.Termination {
		fmt.Fprintln(os.Stderr, "Error:", err)
		log.Fatal(err)
	}
}
