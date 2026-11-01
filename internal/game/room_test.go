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

// Ghosts keep moving through the maze at 4 px a step (never inside a wall),
// and a pellet turns them scared, then flashing at 160 steps, then back.
func TestGhostsAndPellets(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadWorld(pacmanio.Files, a, "ar_1a_nivelo")
	if err != nil {
		t.Fatal(err)
	}
	ghosts := w.Room.All("monster")
	if len(ghosts) == 0 {
		t.Fatal("no ghosts")
	}
	g := ghosts[0]
	x0, y0 := g.X, g.Y
	for range 200 {
		w.Step(held{})
		if !w.Room.PlaceFree(g, g.X, g.Y) {
			t.Fatalf("ghost inside a wall at %v,%v", g.X, g.Y)
		}
	}
	if g.X == x0 && g.Y == y0 {
		t.Error("ghost didn't move")
	}
	p := w.Room.All("pacman")[0]
	pil := w.Room.Create(a.Object("pil"), p.X, p.Y)
	w.Step(held{})
	if !pil.Destroyed || g.Obj.Name != "scared" || g.Alarm[0] <= 0 {
		t.Fatalf("pellet: destroyed %v, ghost %s alarm %d", pil.Destroyed, g.Obj.Name, g.Alarm[0])
	}
	for range 160 {
		w.Step(held{})
	}
	if g.Obj.Name != "intermonster" {
		t.Fatalf("after 160 steps: %s", g.Obj.Name)
	}
	for range 40 {
		w.Step(held{})
	}
	if g.Obj.Name != "monster" {
		t.Errorf("after 40 more: %s", g.Obj.Name)
	}
}

// The fruit in a level: hidden at first, at its place for 200 steps from
// step 500, eaten for 200 points; an extra life adds one.
func TestFruitAndExtraLife(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadWorld(pacmanio.Files, a, "pu_1a_nivelo")
	if err != nil {
		t.Fatal(err)
	}
	f := w.Room.Create(a.Object("bonus"), 64, 64)
	if f.X != -1000 {
		t.Fatalf("not hidden at start: %v", f.X)
	}
	for range 500 {
		w.Step(held{})
	}
	if f.X != 64 || f.Y != 64 {
		t.Fatalf("not back at 500: %v,%v", f.X, f.Y)
	}
	for range 200 {
		w.Step(held{})
	}
	if f.X != -1000 {
		t.Fatalf("not hidden again at 700: %v", f.X)
	}
	p := w.Room.All("pacman")[0]
	f.X, f.Y = p.X, p.Y
	score, lives := w.Room.G.Score, w.Room.G.Lives
	w.Room.Create(a.Object("vivo"), p.X, p.Y)
	w.Step(held{})
	if w.Room.G.Score-score < 200 || f.X != -1000 {
		t.Errorf("fruit not eaten: score +%d, x %v", w.Room.G.Score-score, f.X)
	}
	if w.Room.G.Lives != lives+1 || len(w.Room.All("vivo")) != 0 {
		t.Errorf("extra life: lives %d (was %d)", w.Room.G.Lives, lives)
	}
}

// In every level the regulilo comes after the fruits in the room, so a
// fruit's create event sees none: the level's fruits show from the start
// until step 700 (alarm 0 at 500, alarm 1 200 later hides them).
func TestLevelFruitShowsFirst700Steps(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadWorld(pacmanio.Files, a, "pu_1a_nivelo")
	if err != nil {
		t.Fatal(err)
	}
	f := w.Room.All("bonus")[0]
	x := f.X
	if x < 0 {
		t.Fatal("hidden at the start")
	}
	for range 699 {
		w.Step(held{})
	}
	if f.X != x {
		t.Fatalf("moved before 700: %v", f.X)
	}
	w.Step(held{})
	if f.X != -1000 {
		t.Errorf("not hidden at 700: %v", f.X)
	}
}
