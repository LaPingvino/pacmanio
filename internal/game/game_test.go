package game

import (
	"testing"

	"pacmanio"
	"pacmanio/internal/gm"
)

// The rooms in order: clearing komenco's dot goes on after the 2 s pause;
// the region follows the *_bildo room; Page Down (DEBUG) skips a level with
// score −10000; no more lives goes to fino, which keeps the region.
func TestGameRooms(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(pacmanio.Files, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	clear := func() {
		for _, p := range g.World.Room.All("point") {
			p.Destroyed = true
		}
	}
	step := func(n int, in gm.Input) {
		for range n {
			if err := g.Step(in); err != nil {
				t.Fatal(err)
			}
		}
	}
	if RoomOrder[g.Room] != "komenco" || g.G.Lives != 3 {
		t.Fatalf("start: %s, lives %d", RoomOrder[g.Room], g.G.Lives)
	}
	clear()
	step(1, held{})
	if RoomOrder[g.Room] != "komenco" {
		t.Fatal("left before the pause")
	}
	step(61, held{}) // 2000 ms at room speed 30
	if RoomOrder[g.Room] != "pu_bildo" || g.Region != 2 {
		t.Fatalf("after komenco: %s, region %d", RoomOrder[g.Room], g.Region)
	}
	clear()
	step(45, held{})
	if RoomOrder[g.Room] != "pu_1a_nivelo" {
		t.Fatalf("after pu_bildo: %s", RoomOrder[g.Room])
	}
	step(1, held{gm.KeyPgDn: true})
	if RoomOrder[g.Room] != "pu_2a_nivelo" || g.G.Score != -10000 {
		t.Fatalf("Page Down: %s, score %d", RoomOrder[g.Room], g.G.Score)
	}
	g.G.AddLives(-3)
	step(1, held{})
	if RoomOrder[g.Room] != "fino" || g.Region != 2 {
		t.Errorf("no more lives: %s, region %d", RoomOrder[g.Room], g.Region)
	}
}

// A region title room: the background scrolls right first, then turns
// down, left, up, each after 32, 64 or 96 steps.
func TestRegionTitleScroll(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	w, err := LoadWorld(pacmanio.Files, a, "pu_bildo")
	if err != nil {
		t.Fatal(err)
	}
	var dirs [][2]float64
	last := [2]float64{9, 9}
	for range 500 {
		w.Step(held{})
		d := [2]float64{w.Room.G.Vars["background_hspeed"], w.Room.G.Vars["background_vspeed"]}
		if d != last {
			dirs = append(dirs, d)
			last = d
		}
	}
	want := [][2]float64{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 0}}
	if len(dirs) < len(want) {
		t.Fatalf("directions %v", dirs)
	}
	for i, d := range want {
		if dirs[i] != d {
			t.Fatalf("directions %v, want %v…", dirs, want)
		}
	}
	if w.BgX == 0 && w.BgY == 0 {
		t.Error("the background didn't move")
	}
}

type endKeys struct {
	held
	enter bool
	chars []rune
	any   bool
}

func (e endKeys) Pressed(k gm.Key) bool { return k == gm.KeyEnter && e.enter }
func (e endKeys) Typed() []rune         { return e.chars }
func (e endKeys) AnyPressed() bool      { return e.any }

// The end room: a new highscore asks the name ("Pacman" filled in, Backspace
// edits), Enter puts it in the table; then any key goes back to komenco.
func TestEndRoomHighscore(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(pacmanio.Files, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	g.G.Score = 1230
	for i, n := range RoomOrder {
		if n == "fino" {
			if err := g.enter(i); err != nil {
				t.Fatal(err)
			}
		}
	}
	if !g.asking {
		t.Fatal("a 1230 score should be asked for a name")
	}
	steps := []endKeys{{chars: []rune("\b\b\b\b\b\bJoop")}, {}, {enter: true}}
	for _, k := range steps {
		if err := g.Step(k); err != nil {
			t.Fatal(err)
		}
	}
	if g.asking || g.Highscores[0] != (Highscore{"Joop", 1230}) {
		t.Fatalf("table %v, asking %v", g.Highscores[0], g.asking)
	}
	g.Step(endKeys{any: true})
	if RoomOrder[g.Room] != "komenco" {
		t.Errorf("any key: %s", RoomOrder[g.Room])
	}
}

// Pacman meeting a ghost: a monster costs a life (Pacman and every ghost back
// to their starts); a scared ghost is eaten — +100, back to its start as a
// monster, no life lost (Joop 2026-11-02: "Pacman can eat the blues"), and so
// is a flashing (intermonster) one.
func TestPacmanMeetsGhosts(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"monster", "scared", "intermonster"} {
		w, err := LoadWorld(pacmanio.Files, a, "pu_1a_nivelo")
		if err != nil {
			t.Fatal(err)
		}
		pac := w.Room.All("pacman")[0]
		ghosts := w.Room.All("monster")
		if len(ghosts) == 0 {
			t.Fatal("no ghosts")
		}
		g := ghosts[0]
		if kind != "monster" {
			w.Room.Change(g, a.Object(kind), false)
		}
		pac.X, pac.HSpeed, pac.VSpeed = pac.X+32, 0, 0 // away from his start
		other := ghosts[len(ghosts)-1]
		other.X += 64
		g.X, g.Y, g.HSpeed, g.VSpeed = pac.X, pac.Y, 0, 0
		lives, score := w.Room.G.Lives, w.Room.G.Score
		w.Step(held{})
		switch kind {
		case "monster":
			if w.Room.G.Lives != lives-1 || pac.X != pac.XStart || pac.Y != pac.YStart || other.X != other.XStart {
				t.Errorf("monster: lives %d → %d, pacman (%v,%v) start (%v,%v), other ghost x %v start %v", lives, w.Room.G.Lives, pac.X, pac.Y, pac.XStart, pac.YStart, other.X, other.XStart)
			}
		default:
			eaten := w.Room.All("monster")
			found := false
			for _, e := range eaten {
				if e.X == e.XStart && e.Y == e.YStart && !e.Obj.IsA("scared") && e.Obj.Name == "monster" {
					found = true
				}
			}
			if w.Room.G.Lives != lives || w.Room.G.Score != score+100 || !found {
				t.Errorf("%s: lives %d → %d, score %d → %d, eaten ghost back as a monster %v", kind, lives, w.Room.G.Lives, score, w.Room.G.Score, found)
			}
		}
	}
}
