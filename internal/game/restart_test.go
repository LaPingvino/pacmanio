package game

import (
	"testing"

	"pacmanio"
)

// The pause menu's restart: the level again, with the score and lives it
// began with; not offered outside levels.
func TestRestartLevel(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(pacmanio.Files, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	if g.InLevel() {
		t.Fatal("komenco counts as a level")
	}
	if err := g.StartAt("pu_bildo"); err != nil || g.InLevel() {
		t.Fatalf("pu_bildo: err %v, InLevel %v", err, g.InLevel())
	}
	if err := g.EnterRoom("pu_2a_nivelo"); err != nil || !g.InLevel() {
		t.Fatalf("pu_2a_nivelo: err %v, InLevel %v", err, g.InLevel())
	}
	room, score, lives, old := g.Room, g.G.Score, g.G.Lives, g.World
	g.G.Score, g.G.Lives = score+300, lives-1
	if err := g.RestartLevel(); err != nil {
		t.Fatal(err)
	}
	if g.World == old || g.Room != room || g.G.Score != score || g.G.Lives != lives {
		t.Fatalf("after restart: new world %v, room %d (want %d), score %d (want %d), lives %d (want %d)",
			g.World != old, g.Room, room, g.G.Score, score, g.G.Lives, lives)
	}
}
