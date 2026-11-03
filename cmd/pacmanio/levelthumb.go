package main

import (
	"image"
	"image/color"
	"image/draw"

	"pacmanio"
	"pacmanio/internal/level"
)

// The level pick's pictures (docs/UI.md §3: "each with a small picture of the
// level"): a level's maze, small — the region's background (PACMANIO-RULES §7:
// the City light grey, the Water blue, Free Space dark grey), the walls, the
// dots and coins, the fruit, Pacman and the ghosts in their colours.

var thumbColours = map[string]color.RGBA{
	"konstruu_muron": {0x30, 0x40, 0x90, 255}, "muro": {0x30, 0x40, 0x90, 255},
	"konstruu_murox": {0x50, 0x60, 0x80, 255}, "murox": {0x50, 0x60, 0x80, 255},
	"pil": {0xff, 0xe0, 0x2a, 255}, "monster": {0x2a, 0xd0, 0x4a, 255}, "bonus": {0xe0, 0x2a, 0x3a, 255},
	"pacman": {0xff, 0xf0, 0x00, 255}, "vivo": {0x6a, 0xe8, 0xff, 255},
}

// wallColours keep the walls apart from each region's background.
var wallColours = map[string]color.RGBA{"pu": {0x30, 0x40, 0x90, 255}, "ar": {0xa8, 0xd0, 0xff, 255}, "ls": {0x9a, 0x9a, 0xb4, 255}}

var regionColours = map[string]color.RGBA{"pu": {0xc8, 0xc8, 0xc8, 255}, "ar": {0x20, 0x48, 0xb0, 255}, "ls": {0x40, 0x40, 0x48, 255}}

// levelThumb draws room name's level into a w×h picture (its first screen,
// fitted and centred).
func levelThumb(name string, w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.RGBA{0x05, 0x08, 0x20, 255}), image.Point{}, draw.Src)
	data, err := pacmanio.Files.ReadFile("levels/" + name + ".txt")
	if err != nil {
		return img
	}
	l, err := level.Parse(string(data))
	if err != nil || l.W <= 0 || l.H <= 0 {
		return img
	}
	vw, vh := min(l.W, 640), min(l.H, 480) // no views: the first screen is what shows
	k := min(float64(w)/float64(vw), float64(h)/float64(vh))
	ox, oy := (float64(w)-float64(vw)*k)/2, (float64(h)-float64(vh)*k)/2
	bg, ok := regionColours[name[:min(2, len(name))]]
	if !ok {
		bg = color.RGBA{0, 0, 0, 255}
	}
	draw.Draw(img, image.Rect(int(ox), int(oy), int(ox+float64(vw)*k), int(oy+float64(vh)*k)), image.NewUniform(bg), image.Point{}, draw.Src)
	cell := max(1, int(32*k+0.5))
	dot := max(1, cell/3)
	for _, in := range l.Instances {
		if in.X >= vw || in.Y >= vh {
			continue
		}
		x, y := int(ox+float64(in.X)*k), int(oy+float64(in.Y)*k)
		if in.Object == "point" { // a dot: small, in the middle of its square
			x, y = x+(cell-dot)/2, y+(cell-dot)/2
			draw.Draw(img, image.Rect(x, y, x+dot, y+dot), image.NewUniform(color.RGBA{0xff, 0xf4, 0xc0, 255}), image.Point{}, draw.Src)
			continue
		}
		if c, ok := wallColours[name[:min(2, len(name))]]; ok && (in.Object == "konstruu_muron" || in.Object == "muro") {
			draw.Draw(img, image.Rect(x, y, x+cell, y+cell), image.NewUniform(c), image.Point{}, draw.Src)
			continue
		}
		if c, ok := thumbColours[in.Object]; ok {
			draw.Draw(img, image.Rect(x, y, x+cell, y+cell), image.NewUniform(c), image.Point{}, draw.Src)
		}
	}
	return img
}
