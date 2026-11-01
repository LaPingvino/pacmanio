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
