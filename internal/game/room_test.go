package game

import (
	"testing"

	"pacmanio"
	"pacmanio/internal/gm"
)

// muro's room-start table (pacmanio/objects/muro.gml): every case.
func TestWallFrames(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	cases := map[[4]bool]int{
		{}: 0, {false, false, false, true}: 2, {false, false, true, false}: 4, {false, false, true, true}: 13,
		{false, true, false, false}: 1, {false, true, false, true}: 14, {false, true, true, false}: 10, {false, true, true, true}: 9,
		{true, false, false, false}: 3, {true, false, false, true}: 12, {true, false, true, false}: 15, {true, false, true, true}: 8,
		{true, true, false, false}: 11, {true, true, false, true}: 7, {true, true, true, false}: 6, {true, true, true, true}: 5,
	}
	for n, want := range cases {
		w := &World{Assets: a, Region: 1, Room: &gm.Room{W: 640, H: 480, G: &gm.Globals{Vars: map[string]float64{}}}}
		wall := w.Room.Create(a.Object("muro"), 64, 64)
		for i, d := range [4][2]float64{{-32, 0}, {0, -32}, {32, 0}, {0, 32}} {
			if n[i] {
				w.Room.Create(a.Object("murox"), 64+d[0], 64+d[1])
			}
		}
		w.autotile()
		if got := int(wall.ImageIndex); got != want {
			t.Errorf("neighbours %v: frame %d, want %d", n, got, want)
		}
	}
}

type held map[gm.Key]bool

func (h held) Held(k gm.Key) bool { return h[k] }
func (held) Pressed(gm.Key) bool  { return false }

// pu_1a_nivelo: Pacman starts at (320, 288); left of him are dots every 64 px.
func TestPacmanEatsDotsAndStopsAtWalls(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadWorld(pacmanio.Files, a, "pu_1a_nivelo")
	if err != nil {
		t.Fatal(err)
	}
	if n := len(w.Room.All("konstruu_muron")); n != 0 {
		t.Fatalf("%d placeholders left (should have become walls)", n)
	}
	p := w.Room.All("pacman")[0]
	dots := len(w.Room.All("point"))
	for range 100 {
		w.Step(held{gm.KeyLeft: true})
	}
	if w.Room.G.Score == 0 || len(w.Room.All("point")) >= dots {
		t.Errorf("no dots eaten going left: score %d", w.Room.G.Score)
	}
	if p.HSpeed != 0 || p.Sprite != "pac_stand" {
		t.Errorf("not stopped at the wall: hspeed %v sprite %s at %v,%v", p.HSpeed, p.Sprite, p.X, p.Y)
	}
}
