package main

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"pacmanio/internal/game"
	"pacmanio/internal/store"
)

// The title screen (docs/UI.md §3, Joop: "make it new"): the game's own art and
// a menu in Noto Sans, instead of the 2006 title room. Play starts at the first
// region, Continue at the furthest room reached ("resume wherever you got, like
// with Pits"), Highscores shows the table.

type titleScreen struct {
	on       bool
	menu     game.Menu
	ids      []string
	tick     int
	px, py   int
	progress int    // the furthest room reached (game.RoomOrder index)
	info     string // an open info page: "help" (Kiel ludi), "about" (Pri), "" none
	page     string // the open page: "" the blocks, "settings", "levels" (Daŭrigi's level pick)
	rooms    []int  // the levels page's rooms (RoomOrder indices)
}

func (a *app) openTitle() {
	a.title.on = true
	a.buildTitleMenu(0)
	a.game.G.Play("@music:finsono") // the 2006 title and end tune
	ebiten.SetTPS(30)
	ebiten.SetWindowTitle("Pacmanio")
}

// buildTitleMenu builds the open page of the title (U7, as Stelo; Joop
// 2026-11-04: "block buttons"): the main page as blocks — Ludi, Daŭrigi,
// Plej bonaj, Kiel ludi, Pri, Agordoj (Fini on the desktop); Agordoj the
// settings as a list; Daŭrigi the reached levels as blocks, a region a row.
func (a *app) buildTitleMenu(sel int) {
	a.title.progress = store.Progress()
	var items []game.MenuItem
	var ids []string
	add := func(id, label, value string) {
		items = append(items, game.MenuItem{Label: label, Value: value})
		ids = append(ids, id)
	}
	switch a.title.page {
	case "settings":
		onOff := func(off bool) string {
			if off {
				return "malŝaltita"
			}
			return "ŝaltita"
		}
		add("sound", "Sono", onOff(soundMuted()))
		add("music", "Muziko", onOff(theSound != nil && theSound.MusicOff()))
		add("back", "Reen", "")
		a.title.ids = ids
		a.title.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 230, LineH: 32, W: 330, Size: 18}
		return
	case "levels":
		a.title.rooms = a.title.rooms[:0]
		region := map[string]string{"pu": "Urba", "ar": "Akva", "ls": "Libera"}
		for i := 1; i <= a.title.progress && i < len(game.RoomOrder); i++ {
			name := game.RoomOrder[i]
			r, ok := region[name[:min(2, len(name))]]
			if !ok || strings.HasSuffix(name, "_bildo") {
				continue // the title and end rooms, the region title rooms
			}
			label := "Fina"
			if len(name) > 4 && name[4] == 'a' {
				label = name[3:4]
			}
			add("level", label, r)
			a.title.rooms = append(a.title.rooms, i)
		}
		if sel < 0 || sel >= len(items) {
			sel = len(items) - 1 // the furthest
		}
		a.title.ids = ids
		a.title.menu = game.Menu{Items: items, Sel: sel, X: game.ScreenW / 2, Y: 214, Size: 20, Cols: 5, CellW: 96, CellH: 52, Gap: 10}
		return
	}
	add("play", "Ludi", "")
	if p := a.title.progress; p > 1 && p < len(game.RoomOrder)-1 {
		add("continue", "Daŭrigi", game.RoomTitle(p))
	}
	add("scores", "Plej bonaj", "")
	add("help", "Kiel ludi", "")
	add("about", "Pri", "")
	add("settings", "Agordoj", "")
	if runtime.GOOS != "js" {
		add("quit", "Fini", "")
	}
	cols, w := 3, 172.0
	if len(items) > 6 {
		cols, w = 4, 142
	}
	a.title.ids = ids
	a.title.menu = game.Menu{Items: items, Sel: min(sel, len(items)-1), X: game.ScreenW / 2, Y: 222, Size: 19, Cols: cols, CellW: w, CellH: 62, Gap: 12}
}

// openPage shows a page of the title ("" the main one).
func (a *app) openPage(page string, sel int) {
	a.title.page = page
	a.buildTitleMenu(sel)
}

func (a *app) updateTitle() {
	a.title.tick++
	in := game.MenuInput{
		Up: inpututil.IsKeyJustPressed(ebiten.KeyArrowUp), Down: inpututil.IsKeyJustPressed(ebiten.KeyArrowDown),
		Left: inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft), Right: inpututil.IsKeyJustPressed(ebiten.KeyArrowRight),
		Enter: inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
	}
	x, y := toFrame(ebiten.CursorPosition())
	in.Moved = x != a.title.px || y != a.title.py
	a.title.px, a.title.py = x, y
	in.PX, in.PY = x, y
	in.Clicked = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if ids := inpututil.AppendJustPressedTouchIDs(nil); len(ids) > 0 {
		in.PX, in.PY = toFrame(ebiten.TouchPosition(ids[0]))
		in.Clicked = true
	}
	if a.title.info != "" { // an info page: Enter, Esc or a click goes back
		if in.Enter || in.Clicked || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			a.title.info = ""
		}
		return
	}
	if a.title.page != "" && inpututil.IsKeyJustPressed(ebiten.KeyEscape) { // a page: Esc goes back
		a.openPage("", 0)
		return
	}
	res := a.title.menu.Update(in)
	sel := a.title.menu.Sel
	toggle := func(id string) bool {
		switch id {
		case "sound":
			toggleSound()
		case "music":
			toggleMusic()
		default:
			return false
		}
		a.buildTitleMenu(sel)
		return true
	}
	if res.Step != 0 && toggle(a.title.ids[sel]) {
		return
	}
	if res.Chosen < 0 || toggle(a.title.ids[res.Chosen]) {
		return
	}
	switch a.title.ids[res.Chosen] {
	case "play":
		a.startAt("pu_bildo")
	case "continue":
		a.openPage("levels", -1)
	case "level":
		a.startAt(game.RoomOrder[a.title.rooms[res.Chosen]])
	case "settings":
		a.openPage("settings", 0)
	case "back": // to the main page, Agordoj chosen
		a.openPage("", 0)
		for i, id := range a.title.ids {
			if id == "settings" {
				a.title.menu.Sel = i
			}
		}
	case "scores":
		a.startAt("fino")
	case "help", "about":
		a.title.info = a.title.ids[res.Chosen]
	case "quit":
		os.Exit(0)
	}
}

func (a *app) startAt(room string) {
	if err := a.game.StartAt(room); err != nil {
		log.Printf("room %s: %v", room, err)
	}
	a.title.on = false
}

// drawTitle draws the title: the frame (night blue, a row of Joop's houses,
// Pacman chasing ghosts along the dots) and the overlay (the name, the menu).
func (a *app) drawTitle() {
	dst := a.frame
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.RGBA{0x0b, 0x12, 0x3a, 255}), image.Point{}, draw.Src)
	put := func(name string, frame, x, y int) {
		sp := a.game.Assets.Sprites[name]
		if sp == nil || len(sp.Frames) == 0 {
			return
		}
		img := sp.Frames[frame%len(sp.Frames)]
		draw.Draw(dst, img.Bounds().Add(image.Pt(x, y)), img, img.Bounds().Min, draw.Over)
	}
	for x := 0; x < game.ScreenW; x += 32 { // a street of houses (muro_2: left and right neighbours)
		put("muro_2", 15, x, 448)
	}
	t := a.title.tick
	run := (t * 2) % (game.ScreenW + 256) // the chase crosses the screen
	px := run - 128
	for i := range 20 { // the dots ahead of Pacman
		if x := 16 + i*32; x > px+16 {
			put("point_sprite", 0, x, 404)
		}
	}
	put("pac_right", t/4, px, 404)
	put("afraid", t/6, px+96, 404)
	put("afraid", t/6+1, px+144, 404)
	put("coin_sprite", 0, 568, 360)
	put("bonus_sprite", (t/90)%4, 40, 360)
	o := &a.overlay
	if a.title.info != "" { // the page instead of the title's texts (the overlay draws texts over every rect)
		a.drawInfo(o, a.title.info, put)
		return
	}
	shadow := color.RGBA{0x7a, 0x10, 0x10, 255}
	o.Add(game.Text{S: "Pacmanio", X: game.ScreenW/2 + 4, Y: 54, Size: 88, Weight: game.Black, Color: shadow, Align: game.Centre})
	o.Add(game.Text{S: "Pacmanio", X: game.ScreenW / 2, Y: 50, Size: 88, Weight: game.Black, Color: color.RGBA{0xff, 0xe0, 0x2a, 255}, Align: game.Centre})
	o.Add(game.Text{S: "Labirinta ludo · 2006, refarita 2026", X: game.ScreenW / 2, Y: 176, Size: 15, Weight: game.Bold, Color: color.RGBA{0xd8, 0xe4, 0xff, 255}, Align: game.Centre})
	a.title.menu.Draw(o, game.MenuColors{Text: color.RGBA{0xf0, 0xf4, 0xff, 255}, Disabled: color.RGBA{0x80, 0x88, 0xa8, 255},
		Selected: color.RGBA{0x0b, 0x12, 0x3a, 255}, Highlight: color.RGBA{0xff, 0xe0, 0x2a, 235}})
}
