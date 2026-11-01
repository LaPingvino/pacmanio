// Package game is the Pacmanio remake: rooms drawn and run as the recovered
// 2006 game (docs/PACMANIO-RULES.md).
package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"
	"strings"

	"pacmanio/internal/level"
)

// ScreenW, ScreenH is the window (the original's 640×480).
const ScreenW, ScreenH = 640, 480

// Thing is one instance in a room (drawing state only, for now).
type Thing struct {
	Object string
	X, Y   int
	Sprite string
	Frame  int
}

// Room is a loaded room.
type Room struct {
	Level  *level.Level
	Region int // global.murigxu: 1 Akva (muro_1, blue), 2 Urba (muro_2, light grey), 3 Libera Spaco (muro_3, dark grey)
	Things []*Thing
}

// RegionOf is the region a room is played in: the *_bildo room before it sets
// global.murigxu (pu → 2, ar → 1, ls → 3); komenco sets 1.
// ASSUMPTION: fino keeps the region of the last room played; standalone it is 1.
func RegionOf(room string) int {
	switch {
	case strings.HasPrefix(room, "pu_"):
		return 2
	case strings.HasPrefix(room, "ls_"):
		return 3
	}
	return 1
}

// regionColour is preparu's background_color: c_blue, c_ltgray, c_dkgray
// (Game Maker's colour constants).
var regionColour = map[int]color.RGBA{1: {0, 0, 255, 255}, 2: {192, 192, 192, 255}, 3: {128, 128, 128, 255}}

// LoadRoom reads levels/<name>.txt, turns the wall placeholders into walls
// (konstruu_muron → muro, konstruu_murox → murox: action_change_object at
// create) and picks every wall's sprite and frame as muro's room start does.
func LoadRoom(fsys fs.FS, a *Assets, name string) (*Room, error) {
	data, err := fs.ReadFile(fsys, "levels/"+name+".txt")
	if err != nil {
		return nil, err
	}
	l, err := level.Parse(string(data))
	if err != nil {
		return nil, err
	}
	r := &Room{Level: l, Region: RegionOf(name)}
	for _, in := range l.Instances {
		obj := in.Object
		switch obj {
		case "konstruu_muron":
			obj = "muro"
		case "konstruu_murox":
			obj = "murox"
		}
		t := &Thing{Object: obj, X: in.X, Y: in.Y}
		if d := a.Objects[obj]; d != nil {
			t.Sprite = d.Sprite
		}
		r.Things = append(r.Things, t)
	}
	r.autotile()
	return r, nil
}

// wallFrame is muro/murox's room-start table: the frame for which of the
// neighbours (left, up, right, down) hold a solid (muro or murox).
var wallFrame = [16]int{0, 2, 4, 13, 1, 14, 10, 9, 3, 12, 15, 8, 11, 7, 6, 5}

func (r *Room) autotile() {
	solid := map[[2]int]bool{}
	for _, t := range r.Things {
		if t.Object == "muro" || t.Object == "murox" {
			solid[[2]int{t.X, t.Y}] = true
		}
	}
	bit := func(x, y int) int {
		if solid[[2]int{x, y}] {
			return 1
		}
		return 0
	}
	for _, t := range r.Things {
		if t.Object != "muro" && t.Object != "murox" {
			continue
		}
		if t.Object == "muro" { // global.murigxu picks the region's tile set
			t.Sprite = []string{"", "muro_1", "muro_2", "muro_3"}[r.Region]
		}
		idx := bit(t.X-32, t.Y)<<3 | bit(t.X, t.Y-32)<<2 | bit(t.X+32, t.Y)<<1 | bit(t.X, t.Y+32)
		t.Frame = wallFrame[idx]
	}
}

// Draw renders the room's top-left 640×480 (no views, as the original).
func (r *Room) Draw(dst *image.RGBA, a *Assets) {
	bg := regionColour[r.Region]
	draw.Draw(dst, dst.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)
	things := append([]*Thing(nil), r.Things...)
	depth := func(t *Thing) int {
		if d := a.Objects[t.Object]; d != nil {
			return d.Depth
		}
		return 0
	}
	sort.SliceStable(things, func(i, j int) bool { return depth(things[i]) > depth(things[j]) })
	for _, t := range things {
		d := a.Objects[t.Object]
		sp := a.Sprites[t.Sprite]
		if d == nil || !d.Visible || sp == nil || len(sp.Frames) == 0 {
			continue
		}
		img := sp.Frames[t.Frame%len(sp.Frames)]
		rr := img.Bounds().Add(image.Pt(t.X, t.Y))
		if rr.Overlaps(dst.Bounds()) {
			draw.Draw(dst, rr, img, image.Point{}, draw.Over)
		}
	}
}
