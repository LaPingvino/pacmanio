package game

import (
	"image"
	"testing"

	"pacmanio"
	"pacmanio/internal/gm"
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

// The level-done message: shown in the pause after the last dot, not in a
// death's pause or while playing.
func TestLevelDone(t *testing.T) {
	a, err := LoadAssets(pacmanio.Files)
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGame(pacmanio.Files, a, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.EnterRoom("pu_1a_nivelo"); err != nil {
		t.Fatal(err)
	}
	w := g.World
	if w.LevelDone() {
		t.Fatal("done at the start")
	}
	for _, in := range w.Room.All("point") {
		in.Destroyed = true
	}
	for i := 0; i < 3 && w.Room.Paused == nil; i++ {
		if err := g.Step(noInput{}); err != nil {
			t.Fatal(err)
		}
	}
	if w.Room.Paused == nil || !w.LevelDone() {
		t.Fatalf("paused %v, LevelDone %v", w.Room.Paused != nil, w.LevelDone())
	}
	w.Overlay = &Overlay{}
	w.drawBar(image.NewRGBA(image.Rect(0, 0, ScreenW, ScreenH)), 0, 3)
	found := false
	for _, tx := range w.Overlay.Texts {
		found = found || tx.S == "Nivelo finita!"
		if tx.S == w.Label && tx.Size > 20 {
			t.Error("the start banner shows during the level-done pause")
		}
	}
	if !found {
		t.Error("no 'Nivelo finita!' in the overlay")
	}
}

type noInput struct{}

func (noInput) Held(gm.Key) bool    { return false }
func (noInput) Pressed(gm.Key) bool { return false }
