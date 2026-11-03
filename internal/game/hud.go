package game

import (
	xdraw "golang.org/x/image/draw"
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Fonts stand in for the 2006 Windows fonts: Liberation Sans Bold for
// "MS Sans Serif" 9 bold, Liberation Mono Bold for "Courier New" 18 bold
// (both OFL, fonts/LICENSE-Liberation).
type Fonts struct{ Small, Lives, Region, Table font.Face }

// LoadFonts reads the fonts at Game Maker's sizes (points at 96 DPI).
func LoadFonts(fsys fs.FS) (*Fonts, error) {
	face := func(name string, size float64) (font.Face, error) {
		data, err := fs.ReadFile(fsys, "fonts/"+name)
		if err != nil {
			return nil, err
		}
		f, err := opentype.Parse(data)
		if err != nil {
			return nil, err
		}
		return opentype.NewFace(f, &opentype.FaceOptions{Size: size, DPI: 96, Hinting: font.HintingFull})
	}
	s, err := face("LiberationSans-Bold.ttf", 9)
	if err != nil {
		return nil, err
	}
	l, err := face("LiberationMono-Bold.ttf", 18)
	if err != nil {
		return nil, err
	}
	rg, err := face("LiberationSans-BoldItalic.ttf", 36) // Trebuchet MS 36 bold italic
	if err != nil {
		return nil, err
	}
	tb, err := face("LiberationMono-Bold.ttf", 14) // Courier New 14 bold
	if err != nil {
		return nil, err
	}
	return &Fonts{Small: s, Lives: l, Region: rg, Table: tb}, nil
}

// textColour is global.kontrauxkoloro, set by the region's *_bildo room:
// c_yellow (Akva), c_black (Urba), c_ltgray (Libera Spaco).
var textColour = map[int]color.RGBA{1: {255, 255, 0, 255}, 2: {0, 0, 0, 255}, 3: {192, 192, 192, 255}}

// drawGMText draws Game Maker text: '#' starts a new line.
func drawGMText(dst *image.RGBA, f font.Face, x, y int, s string, c color.Color) {
	m := f.Metrics()
	for _, line := range strings.Split(s, "#") {
		d := font.Drawer{Dst: dst, Src: image.NewUniform(c), Face: f, Dot: fixed.P(x, y+m.Ascent.Ceil())}
		d.DrawString(line)
		y += m.Height.Ceil()
	}
}

// drawHUD is regulilo's draw event: four boxes in the background colour in
// the view's corners (no views: the room's top-left 640×480), with the room
// caption, the score, the programmer and the lives.
// ASSUMPTION: action_draw_rectangle's corners are inclusive (97×65 boxes).
// hudText is drawGMText for the HUD: into the overlay (sharp, Noto Sans Bold
// at the size of MS Sans Serif 9) when there is one, else into the frame.
func (w *World) hudText(dst *image.RGBA, f font.Face, size float64, x, y int, s string, c color.Color) {
	if w.Overlay == nil {
		drawGMText(dst, f, x, y, s, c)
		return
	}
	for i, line := range strings.Split(s, "#") {
		w.Overlay.Add(Text{S: line, X: float64(x), Y: float64(y) + float64(i)*size*1.3, Size: size, Weight: Bold, Color: c})
	}
}

func (w *World) drawHUD(dst *image.RGBA) {
	f := w.Fonts
	if f == nil {
		return
	}
	bg, fg := regionColour[w.Region], textColour[w.Region]
	box := func(x1, y1, x2, y2 int) {
		draw.Draw(dst, image.Rect(x1, y1, x2+1, y2+1), image.NewUniform(bg), image.Point{}, draw.Src)
	}
	g := w.Room.ShownG()  // during a pause, the score and lives of the moment it began
	if w.Overlay != nil { // the new UI's slim bar (docs/UI.md §5)
		w.drawBar(dst, g.Score, g.Lives)
		return
	}
	box(0, 0, 96, 64)
	w.hudText(dst, f.Small, 11, 3, 3, "P A C M A N I O#"+w.Level.Caption, fg)
	box(640-96, 0, 640, 64)
	w.hudText(dst, f.Small, 11, 640-93, 3, "Poentoj:#"+strconv.Itoa(g.Score), fg)
	box(0, 480-64, 96, 480)
	w.hudText(dst, f.Small, 11, 3, 480-61, "Programisto:##Joop KIEFTE", fg)
	box(640-96, 480-64, 640, 480)
	w.hudText(dst, f.Small, 11, 640-93, 480-61, "Restaj vivoj:", fg)
	icon := func(x int) {
		if sp := w.Assets.Sprites["pac_stand"]; sp != nil && len(sp.Frames) > 0 {
			img := sp.Frames[0]
			draw.Draw(dst, img.Bounds().Add(image.Pt(x, 480-32)), img, image.Point{}, draw.Over)
		}
	}
	if g.Lives <= 3 { // !action_if_life(3, larger than): one icon per life
		for i := range max(g.Lives, 0) {
			icon(640 - 96 + 32*i)
		}
	} else { // one icon and the number (Courier New 18 bold, yellow)
		icon(640 - 96)
		w.hudText(dst, f.Lives, 22, 640-64, 480-32, strconv.Itoa(g.Lives), color.RGBA{255, 255, 0, 255})
	}
}

// drawBar is the new UI's HUD (docs/UI.md §5, a remake change): one slim bar
// along the top — the score, the level's title and caption, the lives as
// icons — instead of the four 2006 corner boxes; and the level's title large
// with its caption for its first 40 steps (two seconds at 20 steps a second).
func (w *World) drawBar(dst *image.RGBA, score, lives int) {
	const h = 22
	draw.Draw(dst, image.Rect(0, 0, ScreenW, h), image.NewUniform(color.RGBA{0x0b, 0x12, 0x3a, 190}), image.Point{}, draw.Over)
	text, gold := color.RGBA{0xf0, 0xf4, 0xff, 255}, color.RGBA{0xff, 0xe0, 0x2a, 255}
	w.Overlay.Add(Text{S: "Poentoj  " + strconv.Itoa(score), X: 10, Y: 3, Size: 13, Weight: Bold, Color: text})
	mid := w.Label
	if mid == "" {
		mid = w.Level.Caption
	} else if w.Level.Caption != "" {
		mid += " — " + w.Level.Caption
	}
	w.Overlay.Add(Text{S: mid, X: ScreenW / 2, Y: 3, Size: 13, Weight: Bold, Color: text, Align: Centre})
	if sp := w.Assets.Sprites["pac_stand"]; sp != nil && len(sp.Frames) > 0 {
		img := sp.Frames[0]
		const iw = 18 // the 32 px Pacman, scaled into the bar
		n := max(lives, 0)
		shown := min(n, 5)
		if n > shown {
			shown = 1
		}
		x := ScreenW - 10 - shown*(iw+3)
		if n > 5 {
			w.Overlay.Add(Text{S: "×" + strconv.Itoa(n), X: float64(x - 4), Y: 3, Size: 13, Weight: Bold, Color: gold, Align: Right})
		}
		for i := range shown {
			r := image.Rect(x+i*(iw+3), (h-iw)/2, x+i*(iw+3)+iw, (h-iw)/2+iw)
			xdraw.NearestNeighbor.Scale(dst, r, img, img.Bounds(), xdraw.Over, nil)
		}
	}
	if w.LevelDone() && !w.HideBanner { // within the 2006 pause (2 s), so the timing doesn't change
		w.Overlay.Add(Text{S: "Nivelo finita!", X: ScreenW/2 + 3, Y: ScreenH/2 - 27, Size: 44, Weight: Black, Color: color.NRGBA{0, 0, 0, 255}, Align: Centre})
		w.Overlay.Add(Text{S: "Nivelo finita!", X: ScreenW / 2, Y: ScreenH/2 - 30, Size: 44, Weight: Black, Color: gold, Align: Centre})
		return
	}
	if st := w.Steps; w.Label != "" && st < 40 && !w.HideBanner {
		a := uint8(255)
		if st > 30 { // fading out over the last half second
			a = uint8(255 * (40 - st) / 10)
		}
		w.Overlay.Add(Text{S: w.Label, X: ScreenW/2 + 3, Y: ScreenH/2 - 47, Size: 44, Weight: Black, Color: color.NRGBA{0, 0, 0, a}, Align: Centre})
		w.Overlay.Add(Text{S: w.Label, X: ScreenW / 2, Y: ScreenH/2 - 50, Size: 44, Weight: Black, Color: color.NRGBA{0xff, 0xe0, 0x2a, a}, Align: Centre})
		if w.Level.Caption != "" {
			w.Overlay.Add(Text{S: w.Level.Caption, X: ScreenW/2 + 2, Y: ScreenH/2 + 12, Size: 22, Weight: Bold, Color: color.NRGBA{0, 0, 0, a}, Align: Centre})
			w.Overlay.Add(Text{S: w.Level.Caption, X: ScreenW / 2, Y: ScreenH/2 + 10, Size: 22, Weight: Bold, Color: color.NRGBA{0xf0, 0xf4, 0xff, a}, Align: Centre})
		}
	}
}

// LevelDone reports the pause after a level's last dot (pacman's step: no
// dots left → sound won, 2 s, next room): a pause that began with no dot in
// the room — a death's pause begins with dots left.
func (w *World) LevelDone() bool {
	if w.Label == "" || w.Room.Paused == nil {
		return false
	}
	for _, in := range w.Room.Paused {
		if !in.Destroyed && in.Obj.IsA("point") {
			return false
		}
	}
	return true
}
