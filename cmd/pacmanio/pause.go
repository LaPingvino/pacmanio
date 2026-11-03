package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"pacmanio/internal/game"
)

// The pause menu (docs/UI.md §3, step U2b): Esc or P during a level, or the
// pause button on phones. The game does not step while it is open; the frame
// stays, dimmed. (In the 2006 original Esc ended the program.) As Stelo.
type pauseMenu struct {
	on     bool
	menu   game.Menu
	ids    []string
	px, py int
	page   string // "" the pause menu, "settings" its Agordoj
}

// pauseWanted: Esc or P pressed during a level.
func (a *app) pauseWanted() bool {
	if a.game == nil || !a.game.InLevel() {
		return false
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) || a.touch.pauseTapped
}

func (a *app) openPause() {
	a.pause.on, a.pause.page = true, ""
	a.buildPauseMenu(0)
}

func (a *app) buildPauseMenu(sel int) {
	onOff := func(off bool) string {
		if off {
			return "malŝaltita"
		}
		return "ŝaltita"
	}
	items := []game.MenuItem{ // docs/UI.md §3: Daŭrigi · Rekomenci la nivelon · Agordoj · Ĉefmenuo
		{Label: "Daŭrigi"},
		{Label: "Rekomenci la nivelon"},
		{Label: "Agordoj"},
		{Label: "Ĉefmenuo"},
	}
	a.pause.ids = []string{"resume", "restart", "settings", "title"}
	if a.pause.page == "settings" {
		items = []game.MenuItem{{Label: "Sono", Value: onOff(soundMuted())}, {Label: "Muziko", Value: onOff(theSound != nil && theSound.MusicOff())}, {Label: "Reen"}}
		a.pause.ids = []string{"sound", "music", "back"}
	}
	a.pause.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 190, LineH: 32, W: 330, Size: 17}
}

// updatePause runs the open pause menu.
func (a *app) updatePause() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) || a.touch.pauseTapped {
		if a.pause.page != "" { // Agordoj: back to the pause menu
			a.pause.page = ""
			a.buildPauseMenu(2)
			return
		}
		a.pause.on = false
		return
	}
	in := game.MenuInput{
		Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown),
		Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight),
		Enter: inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	x, y := toFrame(ebiten.CursorPosition())
	in.Moved = x != a.pause.px || y != a.pause.py
	a.pause.px, a.pause.py = x, y
	in.PX, in.PY = x, y
	in.Clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		in.PX, in.PY = toFrame(ebiten.TouchPosition(ids[0]))
		in.Clicked = true
	}
	res := a.pause.menu.Update(in)
	sel := a.pause.menu.Sel
	if res.Step != 0 && (a.pause.ids[sel] == "sound" || a.pause.ids[sel] == "music") {
		if a.pause.ids[sel] == "sound" {
			toggleSound()
		} else {
			toggleMusic()
		}
		a.buildPauseMenu(sel)
		return
	}
	if res.Chosen < 0 {
		return
	}
	switch a.pause.ids[res.Chosen] {
	case "resume":
		a.pause.on = false
	case "restart":
		a.pause.on = false
		if err := a.game.RestartLevel(); err != nil {
			log.Print(err)
		}
	case "sound":
		toggleSound()
		a.buildPauseMenu(sel)
	case "music":
		toggleMusic()
		a.buildPauseMenu(sel)
	case "settings":
		a.pause.page = "settings"
		a.buildPauseMenu(0)
	case "back":
		a.pause.page = ""
		a.buildPauseMenu(2)
	case "title":
		a.pause.on = false
		a.openTitle()
	}
}

// drawPause dims the frame and puts the menu over it.
func (a *app) drawPause() {
	o := &a.overlay
	o.AddRect(game.Rect{X: 0, Y: 0, W: game.ScreenW, H: game.ScreenH, Color: color.RGBA{0, 0, 0, 140}})
	o.AddRect(game.Rect{X: game.ScreenW/2 - 200, Y: 72, W: 400, H: 250, Color: color.RGBA{0x0b, 0x12, 0x3a, 230}}) // a panel behind the menu
	title := "Paŭzo"
	if a.pause.page == "settings" {
		title = "Agordoj"
	}
	o.Add(game.Text{S: title, X: game.ScreenW / 2, Y: 110, Size: 44, Weight: game.Black, Color: color.RGBA{0xff, 0xe0, 0x2a, 255}, Align: game.Centre})
	a.pause.menu.Draw(o, game.MenuColors{Text: color.RGBA{0xf0, 0xf4, 0xff, 255}, Disabled: color.RGBA{0x80, 0x88, 0xa8, 255},
		Selected: color.RGBA{0x0b, 0x12, 0x3a, 255}, Highlight: color.RGBA{0xff, 0xe0, 0x2a, 235}})
}
