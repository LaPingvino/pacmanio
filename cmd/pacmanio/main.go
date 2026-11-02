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
	"pacmanio/internal/store"
)

type app struct {
	game   *game.Game
	frame  *image.RGBA
	canvas *ebiten.Image // the frame as an image, scaled onto the screen
	touch  touchState
}

// keys is the keyboard, the mouse and the touch screen as the game's input.
type keys struct {
	t     *touchState
	phone phoneKeys // this frame's characters and Enter from the phone keyboard
}

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

func (k keys) Held(key gm.Key) bool {
	return ebiten.IsKeyPressed(ebitenKey(key)) || (k.t != nil && k.t.held[key])
}
func (k keys) Pressed(key gm.Key) bool {
	return inpututil.IsKeyJustPressed(ebitenKey(key)) || (key == gm.KeyEnter && k.phone.enter)
}
func (k keys) Typed() []rune {
	t := ebiten.AppendInputChars(nil)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		t = append(t, '\b')
	}
	return append(t, k.phone.chars...)
}
func (k keys) AnyPressed() bool {
	return len(inpututil.AppendJustPressedKeys(nil)) > 0 || inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		(k.t != nil && k.t.tapped) || len(k.phone.chars) > 0 || k.phone.enter
}

func (a *app) Update() error {
	a.touch.update()
	if err := a.game.Step(keys{&a.touch, pollPhoneKeys()}); err != nil {
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
	a.canvas.WritePixels(a.frame.Pix)
	a.touch.draw(a.canvas)
	present(screen, a.canvas)
}

// Layout is the screen's real resolution; the frame is scaled onto it (scale.go).
func (a *app) Layout(outW, outH int) (int, int) { return screenSize(outW, outH) }

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
	player := sound.NewPlayer(pacmanio.Files)
	g.G.Sound, g.G.SoundPlaying = player.Play, player.Playing
	g.Highscores = store.LoadHighscores() // kept between games (localStorage / config file)
	g.SaveHighscores = store.SaveHighscores
	a := &app{game: g, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH)), canvas: ebiten.NewImage(game.ScreenW, game.ScreenH)}
	a.touch.seen = showPadsFromURL()
	ebiten.SetWindowTitle("Pacmanio")
	ebiten.SetTPS(g.World.Speed)
	ebiten.SetWindowSize(game.ScreenW, game.ScreenH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}
