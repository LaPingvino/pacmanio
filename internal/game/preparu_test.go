package game

import (
	"testing"

	"pacmanio/internal/gm"
)

// preparu: each level and region room start loops muziko unless it already
// plays — so it carries on from room to room without restarting.
func TestPreparuMusic(t *testing.T) {
	var played []string
	playing := ""
	r := &gm.Room{G: &gm.Globals{Vars: map[string]float64{},
		Sound:        func(n string) { played = append(played, n); playing = n },
		SoundPlaying: func(n string) bool { return "@music:"+n == playing }}}
	preparu(r)
	preparu(r)
	if len(played) != 1 || played[0] != "@music:muziko" {
		t.Fatalf("played %v", played)
	}
}
