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
	if w.Overlay != nil {
		w.drawEndUI(h, asking, name)
		return
	}
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

// drawEndUI is the new UI's highscore screen (docs/UI.md §5, step U5): a
// table in Noto Sans on a panel — rank, name, score — and, while a name is
// asked, the 2006 question over a text field with a caret. The texts are the
// 2006 ones.
func (w *World) drawEndUI(h Highscores, asking bool, name string) {
	o := w.Overlay
	navy, text, gold := color.RGBA{0x0b, 0x12, 0x3a, 235}, color.RGBA{0xf0, 0xf4, 0xff, 255}, color.RGBA{0xff, 0xe0, 0x2a, 255}
	cx := float64(ScreenW) / 2
	if asking {
		o.AddRect(Rect{X: cx - 260, Y: 150, W: 520, H: 150, Color: navy})
		o.Add(Text{S: hsAsk, X: cx, Y: 168, Size: 16, Weight: Bold, Color: gold, Align: Centre})
		o.AddRect(Rect{X: cx - 200, Y: 210, W: 400, H: 44, Color: color.RGBA{0xf0, 0xf4, 0xff, 255}})
		caret := ""
		if w.Steps/15%2 == 0 { // blinking, half a second each way at 30 steps a second
			caret = "|"
		}
		o.Add(Text{S: name + caret, X: cx - 190, Y: 219, Size: 20, Weight: Bold, Color: color.RGBA{0x0b, 0x12, 0x3a, 255}})
		return
	}
	o.AddRect(Rect{X: cx - 200, Y: 40, W: 400, H: 400, Color: navy})
	o.Add(Text{S: "Plej bonaj", X: cx, Y: 54, Size: 30, Weight: Black, Color: gold, Align: Centre})
	for i, e := range h {
		y := 112 + float64(i)*28
		col := text
		if i == 0 {
			col = gold
		}
		o.Add(Text{S: strconv.Itoa(i+1) + ".", X: cx - 160, Y: y, Size: 17, Weight: Bold, Color: col, Align: Right})
		o.Add(Text{S: e.Name, X: cx - 145, Y: y, Size: 17, Weight: Bold, Color: col})
		o.Add(Text{S: strconv.Itoa(e.Score), X: cx + 170, Y: y, Size: 17, Weight: Bold, Color: col, Align: Right})
	}
	o.Add(Text{S: "Premu klavon por ekludi...", X: cx, Y: 405, Size: 15, Weight: Bold, Color: text, Align: Centre})
}
