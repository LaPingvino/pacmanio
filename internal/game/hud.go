package game

import (
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
type Fonts struct{ Small, Lives font.Face }

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
	return &Fonts{Small: s, Lives: l}, nil
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
func (w *World) drawHUD(dst *image.RGBA) {
	f := w.Fonts
	if f == nil {
		return
	}
	bg, fg := regionColour[w.Region], textColour[w.Region]
	box := func(x1, y1, x2, y2 int) {
		draw.Draw(dst, image.Rect(x1, y1, x2+1, y2+1), image.NewUniform(bg), image.Point{}, draw.Src)
	}
	g := w.Room.G
	box(0, 0, 96, 64)
	drawGMText(dst, f.Small, 3, 3, "P A C M A N I O#"+w.Level.Caption, fg)
	box(640-96, 0, 640, 64)
	drawGMText(dst, f.Small, 640-93, 3, "Poentoj:#"+strconv.Itoa(g.Score), fg)
	box(0, 480-64, 96, 480)
	drawGMText(dst, f.Small, 3, 480-61, "Programisto:##Joop KIEFTE", fg)
	box(640-96, 480-64, 640, 480)
	drawGMText(dst, f.Small, 640-93, 480-61, "Restaj vivoj:", fg)
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
		drawGMText(dst, f.Lives, 640-64, 480-32, strconv.Itoa(g.Lives), color.RGBA{255, 255, 0, 255})
	}
}
