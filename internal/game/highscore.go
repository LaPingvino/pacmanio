package game

import (
	"image"
	"image/color"
	"image/draw"
	"sort"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// The end room (pacmanio/objects/finigxu.gml, in fino):
//
//	Create: music off; if highscore_value(10) < score: highscore_add(
//	        get_string("Vi estas poentegisto! Kio fakte estas via nomo?",
//	        "Pacman"), score) — Windows' input dialog; here an in-game text
//	        field (Joop's decision), "Pacman" filled in; sound_loop(finsono).
//	Step:   a mouse button or any key → komenco.
//	Draw:   Courier New 14 bold yellow; draw_highscore(x, y, x+256, y+256);
//	        "Premu klavon por ekludi..." at (x, y+288).

// Highscore is one entry of the table.
type Highscore struct {
	Name  string
	Score int
}

// Highscores is Game Maker's ten-entry table, best first.
// ASSUMPTION [engine]: an empty table holds "<nobody>" with 0 (as Stelo).
type Highscores [10]Highscore

// ClearedHighscores is highscore_clear().
func ClearedHighscores() Highscores {
	var h Highscores
	for i := range h {
		h[i] = Highscore{"<nobody>", 0}
	}
	return h
}

// Add is highscore_add: the entry goes in its place, the last one drops out.
func (h *Highscores) Add(name string, score int) {
	l := append(h[:], Highscore{name, score})
	sort.SliceStable(l, func(i, j int) bool { return l[i].Score > l[j].Score })
	copy(h[:], l)
}

const hsAsk = "Vi estas poentegisto! Kio fakte estas via nomo?"

var hsYellow = image.NewUniform(color.RGBA{255, 255, 0, 255})

// drawEnd is finigxu's draw for the instance at (x, y) (and, while a name
// is asked, the question with the name being typed).
// ASSUMPTION [engine]: draw_highscore puts names at the left edge and scores
// right-aligned at the right edge, the ten rows spread over the height.
func (w *World) drawEnd(dst *image.RGBA, x, y int, h Highscores, asking bool, name string) {
	f := w.Fonts
	if f == nil || f.Table == nil {
		return
	}
	if asking {
		// The question in a box of the room's background colour, centred
		// (the original's get_string was a Windows dialog over the game).
		m := f.Table.Metrics()
		tw := font.MeasureString(f.Table, hsAsk).Ceil()
		lh := m.Height.Ceil()
		bx, by := (ScreenW-tw)/2-8, y-8
		draw.Draw(dst, image.Rect(bx, by, bx+tw+16, by+2*lh+16), image.NewUniform(regionColour[w.Region]), image.Point{}, draw.Src)
		drawGMText(dst, f.Table, bx+8, by+8, hsAsk+"#"+name+"_", hsYellow.C)
		return
	}
	m := f.Table.Metrics()
	x1, y1, x2, y2 := x, y, x+256, y+256
	for i, e := range h {
		yy := y1 + i*(y2-y1)/10 + m.Ascent.Ceil()
		d := font.Drawer{Dst: dst, Src: hsYellow, Face: f.Table, Dot: fixed.P(x1, yy)}
		d.DrawString(e.Name)
		s := strconv.Itoa(e.Score)
		wd := font.MeasureString(f.Table, s).Ceil()
		d = font.Drawer{Dst: dst, Src: hsYellow, Face: f.Table, Dot: fixed.P(x2-wd, yy)}
		d.DrawString(s)
	}
	drawGMText(dst, f.Table, x, y+288, "Premu klavon por ekludi...", hsYellow.C)
}
