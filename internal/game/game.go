package game

import (
	"io/fs"
	"strings"

	"pacmanio/internal/gm"
)

// RoomOrder is project.json's room order.
var RoomOrder = []string{"komenco",
	"pu_bildo", "pu_1a_nivelo", "pu_2a_nivelo", "pu_3a_nivelo", "pu_4a_nivelo", "pu_fina_nivelo",
	"ar_bildo", "ar_1a_nivelo", "ar_2a_nivelo", "ar_3a_nivelo", "ar_4a_nivelo", "ar_fina_nivelo",
	"ls_bildo", "ls_1a_nivelo", "ls_2a_nivelo", "ls_3a_nivelo", "ls_4a_nivelo", "ls_fina_nivelo",
	"fino"}

// Game runs the rooms in order with the game-wide state.
type Game struct {
	FS     fs.FS
	Assets *Assets
	Fonts  *Fonts
	G      *gm.Globals
	World  *World
	Room   int // index in RoomOrder
	Region int // global.murigxu, kept across rooms (fino keeps the last one)
}

// NewGame starts at the first room.
func NewGame(fsys fs.FS, a *Assets, f *Fonts) (*Game, error) {
	g := &Game{FS: fsys, Assets: a, Fonts: f, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}, Region: 1}
	return g, g.enter(0)
}

// enter goes to room i: the *_bildo room's creation code sets the region
// (pu 2, ar 1, ls 3), komenco sets 1, other rooms keep it.
func (g *Game) enter(i int) error {
	name := RoomOrder[i]
	switch {
	case name == "komenco":
		g.Region = 1
	case strings.HasSuffix(name, "_bildo"):
		g.Region = RegionOf(name)
	}
	w, err := LoadWorldWith(g.FS, g.Assets, name, g.G, g.Region)
	if err != nil {
		return err
	}
	w.Fonts = g.Fonts
	g.World, g.Room = w, i
	return nil
}

// Step runs a step and then any room change it asked for (room_goto_next,
// a named room, or "no more lives" → fino).
func (g *Game) Step(in gm.Input) error {
	r := g.World.Room
	r.Step(in)
	if r.Sleep > 0 {
		return nil // the change waits until the pause is over
	}
	to := r.Goto
	r.Goto = ""
	if g.G.NoMoreLives { // pacman's no-more-lives event: room fino
		g.G.NoMoreLives = false
		to = "fino"
	}
	switch to {
	case "":
		return nil
	case "@next":
		if g.Room+1 < len(RoomOrder) {
			return g.enter(g.Room + 1)
		}
		return nil
	}
	for i, n := range RoomOrder {
		if n == to {
			return g.enter(i)
		}
	}
	return nil
}

func init() {
	// regulilo: Page Down (keyboard event 34, every step it is held) — with
	// DEBUG true as released: score = −10000 and the next room.
	objectEvents["regulilo"] = func(a *Assets, o *gm.Object) {
		o.Events.Key = map[gm.Key]gm.Handler{gm.KeyPgDn: func(r *gm.Room, self *gm.Instance) {
			r.G.Score = -10000
			r.Goto = "@next"
		}}
	}
	// komencu (the title room): music off, score 0, lives 3, finsono, and
	// after 450 steps the highscores (fino, transition 14).
	// PARTIAL: a mouse click quits the game (game_end); the transition.
	objectEvents["komencu"] = func(a *Assets, o *gm.Object) {
		o.Events.Create = func(r *gm.Room, self *gm.Instance) {
			r.G.Score, r.G.Lives = 0, 3
			r.G.Play("@music:finsono")
			self.Alarm[0] = 450
		}
		o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) { r.Goto = "fino" }
	}
}
