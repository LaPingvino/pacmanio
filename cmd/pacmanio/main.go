// Command pacmanio is the Pacmanio remake (Ebitengine; desktop and
// WebAssembly). PRIVATE: it ships the original's stock assets.
package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"pacmanio"
	"pacmanio/internal/game"
	"pacmanio/internal/gm"
	"pacmanio/internal/sound"
)

type app struct {
	game  *game.Game
	frame *image.RGBA
}

// keys is the keyboard (and mouse) as the game's input.
type keys struct{}

func ebitenKey(k gm.Key) ebiten.Key {
	switch k {
	case gm.KeyLeft:
		return ebiten.KeyArrowLeft
	case gm.KeyRight:
		return ebiten.KeyArrowRight
	case gm.KeyUp:
		return ebiten.KeyArrowUp
	case gm.KeyDown:
		return ebiten.KeyArrowDown
	case gm.KeySpace:
		return ebiten.KeySpace
	case gm.KeyPgDn:
		return ebiten.KeyPageDown
	case gm.KeyEnter:
		return ebiten.KeyEnter
	}
	return -1
}

func (keys) Held(k gm.Key) bool    { return ebiten.IsKeyPressed(ebitenKey(k)) }
func (keys) Pressed(k gm.Key) bool { return inpututil.IsKeyJustPressed(ebitenKey(k)) }
func (keys) Typed() []rune {
	t := ebiten.AppendInputChars(nil)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		t = append(t, '\b')
	}
	return t
}
func (keys) AnyPressed() bool {
	return len(inpututil.AppendJustPressedKeys(nil)) > 0 || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
}

func (a *app) Update() error {
	if err := a.game.Step(keys{}); err != nil {
		log.Print(err)
	}
	// Each room has its own speed (levels 20, region titles 15, komenco and
	// fino 30 steps per second).
	if s := a.game.World.Speed; s > 0 && ebiten.TPS() != s {
		ebiten.SetTPS(s)
	}
	return nil
}

func (a *app) Draw(screen *ebiten.Image) {
	a.game.Draw(a.frame)
	screen.WritePixels(a.frame.Pix)
}

func (a *app) Layout(int, int) (int, int) { return game.ScreenW, game.ScreenH }

func main() {
	assets, err := game.LoadAssets(pacmanio.Files)
	if err != nil {
		log.Fatal(err)
	}
	fonts, err := game.LoadFonts(pacmanio.Files)
	if err != nil {
		log.Fatal(err)
	}
	g, err := game.NewGame(pacmanio.Files, assets, fonts)
	if err != nil {
		log.Fatal(err)
	}
	g.G.Sound = sound.NewPlayer(pacmanio.Files).Play
	// PARTIAL: the highscore table isn't kept between games yet (store, as Stelo).
	a := &app{game: g, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))}
	ebiten.SetWindowTitle("Pacmanio")
	ebiten.SetTPS(g.World.Speed)
	ebiten.SetWindowSize(game.ScreenW, game.ScreenH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}
