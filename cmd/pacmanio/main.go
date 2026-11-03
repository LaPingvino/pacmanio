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
	pause   pauseMenu
	game    *game.Game
	frame   *image.RGBA
	canvas  *ebiten.Image // the frame as an image, scaled onto the screen
	overlay game.Overlay  // texts drawn sharp over it (text.go)
	touch   touchState
	title   titleScreen
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
	// The two corner buttons by mouse (no touch layout), and M for the sound.
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !a.touch.seen {
		pt := image.Pt(ebiten.CursorPosition())
		switch {
		case pt.In(a.touch.lay.full):
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			return nil
		case pt.In(a.touch.lay.snd):
			toggleSound()
			return nil
		case a.touch.showPause && pt.In(a.touch.lay.paus):
			if a.pause.on {
				a.pause.on = false
			} else {
				a.openPause()
			}
			return nil
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyM) && !a.typing() {
		toggleSound()
	}
	// Full screen: F11, or Alt+Enter (the touch controls have a button).
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
		return nil
	}
	if a.title.on {
		a.touch.update() // the controls appear (and full screen works) on the title too
		a.updateTitle()
		return nil
	}
	a.touch.update()
	a.touch.showPause = a.game.InLevel()
	if a.pause.on {
		a.updatePause()
		return nil
	}
	if a.pauseWanted() {
		a.openPause()
		return nil
	}
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
	if a.title.on {
		a.overlay.Reset()
		a.drawTitle()
		a.canvas.WritePixels(a.frame.Pix)
		present(screen, a.canvas, a.frameArea(screen))
		drawOverlay(screen, &a.overlay)
		a.touch.draw(screen)
		return
	}
	a.overlay.Reset()
	if a.game.World != nil {
		a.game.World.Overlay = &a.overlay // the HUD's text, drawn sharp after the frame (text.go)
		a.game.World.HideBanner = a.pause.on
	}
	a.game.Draw(a.frame)
	if a.pause.on {
		a.drawPause()
	}
	a.canvas.WritePixels(a.frame.Pix)
	present(screen, a.canvas, a.frameArea(screen))
	drawOverlay(screen, &a.overlay)
	a.touch.draw(screen)
}

// Layout is the screen's real resolution; the frame is scaled onto it (scale.go).
func (a *app) Layout(outW, outH int) (int, int) { return screenSize(outW, outH) }

func main() {
	loadNoto(pacmanio.Files)
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
	theSound = player
	player.SetMuted(store.Muted())
	player.SetMusicOff(store.MusicOff())
	g.G.Sound, g.G.SoundPlaying = player.Play, player.Playing
	g.Highscores = store.LoadHighscores() // kept between games (localStorage / config file)
	g.SaveHighscores = store.SaveHighscores
	a := &app{game: g, frame: image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH)), canvas: ebiten.NewImage(game.ScreenW, game.ScreenH)}
	if r := roomFromURL(); r != "" { // ?room= / PACMANIO_ROOM: start in that room
		if err := g.EnterRoom(r); err != nil {
			log.Printf("room %s: %v", r, err)
		}
	}
	a.touch.seen = showPadsFromURL()
	g.OnRoom = func(i int) { store.Reached(i) } // the furthest room, for Continue
	g.OnTitle = func() { a.openTitle() }        // after the highscores: the new title
	if roomFromURL() == "" {
		a.openTitle() // the new title (docs/UI.md); the 2006 title room no longer shows
	}
	ebiten.SetWindowTitle("Pacmanio")
	ebiten.SetTPS(g.World.Speed)
	ebiten.SetWindowSize(game.ScreenW, game.ScreenH)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(a); err != nil {
		log.Fatal(err)
	}
}

// frameArea is where the frame goes on the screen: all of it, or — once touch
// has been used — what the touch controls leave (touch.go); it also stores the
// layout the touches are read against.
func (a *app) frameArea(screen *ebiten.Image) image.Rectangle {
	b := screen.Bounds()
	if !a.touch.seen { // no touch: the whole screen; the small buttons in the top-left corner
		d := ebiten.Monitor().DeviceScaleFactor()
		k := int(28 * max(d, 1))
		m := k / 4
		h := k * 2 / 3 // in a column, so they stay in the side margin
		a.touch.lay = touchLayout{area: b, full: image.Rect(m, m, m+k, m+h), snd: image.Rect(m, 2*m+h, m+k, 2*m+2*h),
			paus: image.Rect(m, 3*m+2*h, m+k, 3*m+3*h)}
		return b
	}
	a.touch.lay = layoutTouch(b.Dx(), b.Dy())
	return a.touch.lay.area
}

// theSound is the program's sound player.
var theSound *sound.Player

// toggleSound is the sound button and M: all sound off or on, remembered.
func toggleSound() {
	if theSound == nil {
		return
	}
	theSound.SetMuted(!theSound.Muted())
	store.SetMuted(theSound.Muted())
}

func soundMuted() bool { return theSound != nil && theSound.Muted() }

// typing reports the highscore name being typed (M is a letter then).
func (a *app) typing() bool { return a.game != nil && a.game.Asking() }

// toggleMusic is the music setting: the tune alone off or on, remembered.
func toggleMusic() {
	if theSound == nil {
		return
	}
	theSound.SetMusicOff(!theSound.MusicOff())
	store.SetMusicOff(theSound.MusicOff())
}
