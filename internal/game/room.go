// Package game is the Pacmanio remake: rooms run on the Game Maker core
// (internal/gm) as the recovered 2006 game (docs/PACMANIO-RULES.md).
package game

import (
	"image"
	"image/color"
	"image/draw"
	"io/fs"
	"sort"
	"strings"

	"pacmanio/internal/gm"
	"pacmanio/internal/level"
)

// ScreenW, ScreenH is the window (the original's 640×480).
const ScreenW, ScreenH = 640, 480

// World is a running room.
type World struct {
	Assets *Assets
	Level  *level.Level
	Room   *gm.Room
	Region int    // global.murigxu: 1 Akva (muro_1, blue), 2 Urba (muro_2, light grey), 3 Libera Spaco (muro_3, dark grey)
	Speed  int    // room speed (steps per second)
	Fonts  *Fonts // the text's fonts; nil = no text
	// BgX, BgY is the background's offset: it moves by background_hspeed/
	// vspeed each step.
	BgX, BgY float64
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

// regionName is global.regiono, set by the *_bildo rooms' creation code.
var regionName = map[int]string{1: "Akva Regiono", 2: "Urba Regiono", 3: "Libera Spaco"}

// regionColour is preparu's background_color: c_blue, c_ltgray, c_dkgray.
var regionColour = map[int]color.RGBA{1: {0, 0, 255, 255}, 2: {192, 192, 192, 255}, 3: {128, 128, 128, 255}}

// Object returns the engine object for an original object (sprite box,
// solidity, drawing defaults, parent, events), shared by its instances.
func (a *Assets) Object(name string) *gm.Object {
	if o := a.gmObjects[name]; o != nil {
		return o
	}
	if a.gmObjects == nil {
		a.gmObjects = map[string]*gm.Object{}
	}
	o := &gm.Object{Name: name}
	a.gmObjects[name] = o
	if d := a.Objects[name]; d != nil {
		o.Solid, o.Sprite, o.Visible, o.Depth = d.Solid, d.Sprite, d.Visible, d.Depth
		if sp := a.Sprites[d.Sprite]; sp != nil {
			o.Mask = true
			o.BBox = gm.Box{L: sp.BBox[0], R: sp.BBox[1], T: sp.BBox[2], B: sp.BBox[3]}
		}
		if d.Parent != "" {
			o.Parent = a.Object(d.Parent)
		}
	}
	if ev := objectEvents[name]; ev != nil {
		ev(a, o)
	}
	return o
}

// LoadWorld reads levels/<name>.txt and runs the room's start: instances are
// created in order (their create events run), then muro's room start picks
// every wall's tile set and frame.
func LoadWorld(fsys fs.FS, a *Assets, name string) (*World, error) {
	return LoadWorldWith(fsys, a, name, &gm.Globals{Lives: 3, Vars: map[string]float64{}}, RegionOf(name))
}

// LoadWorldWith loads a room with the game's globals (score, lives across
// rooms) in the given region.
func LoadWorldWith(fsys fs.FS, a *Assets, name string, g *gm.Globals, region int) (*World, error) {
	data, err := fs.ReadFile(fsys, "levels/"+name+".txt")
	if err != nil {
		return nil, err
	}
	l, err := level.Parse(string(data))
	if err != nil {
		return nil, err
	}
	w := &World{Assets: a, Level: l, Region: region, Speed: l.Speed,
		Room: &gm.Room{Name: l.Name, W: l.W, H: l.H, G: g}}
	w.Room.G.Vars["room_speed"] = float64(l.Speed)
	// background_hspeed/vspeed belong to the room: they start at 0.
	w.Room.G.Vars["background_hspeed"], w.Room.G.Vars["background_vspeed"] = 0, 0
	for _, in := range l.Instances {
		w.Room.Create(a.Object(in.Object), float64(in.X), float64(in.Y))
	}
	w.autotile()
	return w, nil
}

// wallFrame is muro/murox's room-start table: the frame for which of the
// neighbours (left, up, right, down) hold a solid (place_free false).
var wallFrame = [16]int{0, 2, 4, 13, 1, 14, 10, 9, 3, 12, 15, 8, 11, 7, 6, 5}

func (w *World) autotile() {
	r := w.Room
	for _, t := range r.Instances {
		name := t.Obj.Name
		if name != "muro" && name != "murox" {
			continue
		}
		if name == "muro" { // global.murigxu picks the region's tile set
			t.Sprite = []string{"", "muro_1", "muro_2", "muro_3"}[w.Region]
		}
		bit := func(dx, dy float64) int {
			if r.PlaceFree(t, t.X+dx, t.Y+dy) {
				return 0
			}
			return 1
		}
		idx := bit(-32, 0)<<3 | bit(0, -32)<<2 | bit(32, 0)<<1 | bit(0, 32)
		t.ImageIndex, t.ImageSpeed = float64(wallFrame[idx]), 0 // image_single
	}
}

// Step runs one step.
func (w *World) Step(in gm.Input) {
	w.Room.Step(in)
	w.BgX += w.Room.G.Vars["background_hspeed"]
	w.BgY += w.Room.G.Vars["background_vspeed"]
}

// Draw renders the room's top-left 640×480 (no views, as the original).
func (w *World) Draw(dst *image.RGBA) {
	draw.Draw(dst, dst.Bounds(), image.NewUniform(regionColour[w.Region]), image.Point{}, draw.Src)
	if bg := w.Assets.Backgrounds[w.Level.Background]; bg != nil { // tiled both ways
		bw, bh := bg.Bounds().Dx(), bg.Bounds().Dy()
		ox := ((int(w.BgX) % bw) + bw) % bw
		oy := ((int(w.BgY) % bh) + bh) % bh
		for y := oy - bh; y < dst.Bounds().Dy(); y += bh {
			for x := ox - bw; x < dst.Bounds().Dx(); x += bw {
				draw.Draw(dst, bg.Bounds().Add(image.Pt(x, y)), bg, bg.Bounds().Min, draw.Src)
			}
		}
	}
	insts := append([]*gm.Instance(nil), w.Room.Drawn()...) // during a pause, the moment it began
	sort.SliceStable(insts, func(i, j int) bool { return insts[i].Obj.Depth > insts[j].Obj.Depth })
	for _, t := range insts {
		if t.Obj.Name == "regulilo" && !t.Destroyed { // its draw event (depth −10)
			w.drawHUD(dst)
			continue
		}
		if t.Obj.Name == "regionilo" && !t.Destroyed && w.Fonts != nil {
			// The region's name in white, then in black 1 px down-right.
			name := regionName[w.Region]
			drawGMText(dst, w.Fonts.Region, int(t.X), int(t.Y), name, color.White)
			drawGMText(dst, w.Fonts.Region, int(t.X)+1, int(t.Y)+1, name, color.Black)
			continue
		}
		sp := w.Assets.Sprites[t.Sprite]
		if !t.Obj.Visible || t.Destroyed || sp == nil || len(sp.Frames) == 0 {
			continue
		}
		img := sp.Frames[int(t.ImageIndex)%len(sp.Frames)]
		rr := img.Bounds().Add(image.Pt(int(t.X), int(t.Y)))
		if rr.Overlaps(dst.Bounds()) {
			draw.Draw(dst, rr, img, image.Point{}, draw.Over)
		}
	}
}
