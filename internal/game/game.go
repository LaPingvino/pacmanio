package game

import (
	"image"
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
	// Highscores is the table; SaveHighscores (if set) keeps it.
	Highscores     Highscores
	SaveHighscores func(Highscores)
	asking         bool // fino: the name is being typed (get_string)
	// OnRoom, when set, is told each room entered (its index in RoomOrder):
	// the program keeps the furthest one for Continue.
	OnRoom func(i int)
	// OnTitle, when set, replaces going back to komenco after the highscores:
	// the program's own title screen (docs/UI.md).
	OnTitle func()
}

// AnyKey is keyboard_key != 0 or a mouse button this step.
type AnyKey interface{ AnyPressed() bool }

// Typer gives the characters typed since the last step.
type Typer interface{ Typed() []rune }

// NewGame starts at the first room.
func NewGame(fsys fs.FS, a *Assets, f *Fonts) (*Game, error) {
	g := &Game{FS: fsys, Assets: a, Fonts: f, G: &gm.Globals{Lives: 3, Vars: map[string]float64{}}, Region: 1,
		Highscores: ClearedHighscores()}
	return g, g.enter(0)
}

// enter goes to room i: the *_bildo room's creation code sets the region
// (pu 2, ar 1, ls 3), komenco sets 1, other rooms keep it.
// RoomTitle names a room for the player: a region's name, or the region and
// the level number ("Urba Regiono 2"), or the final level's.
func RoomTitle(i int) string {
	name := RoomOrder[i]
	region := map[string]string{"pu": "Urba Regiono", "ar": "Akva Regiono", "ls": "Libera Spaco"}[name[:2]]
	switch {
	case region == "":
		return name
	case len(name) > 3 && name[3:] == "bildo":
		return region
	case len(name) > 5 && name[3:] == "fina_nivelo":
		return region + " — fina nivelo"
	case len(name) > 4 && name[4] == 'a':
		return region + " " + name[3:4]
	}
	return region
}

// StartAt starts a new game in room name: score 0 and 3 lives, as komencu sets
// them (the title's Play and Continue).
func (g *Game) StartAt(name string) error {
	g.G.Score, g.G.Lives = 0, 3
	return g.EnterRoom(name)
}

// EnterRoom starts in a room by name (a testing aid: ?room= / PACMANIO_ROOM);
// unknown names are ignored.
func (g *Game) EnterRoom(name string) error {
	for i, n := range RoomOrder {
		if n == name {
			return g.enter(i)
		}
	}
	return nil
}

func (g *Game) enter(i int) error {
	if g.OnRoom != nil {
		g.OnRoom(i)
	}
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
	if name == "fino" { // finigxu's create: a new highscore asks the name
		g.asking = g.Highscores[9].Score < g.G.Score
		g.G.KeyboardString = "Pacman"
		g.G.Play("@music:finsono")
	}
	return nil
}

// Step runs a step and then any room change it asked for (room_goto_next,
// a named room, or "no more lives" → fino).
func (g *Game) Step(in gm.Input) error {
	if RoomOrder[g.Room] == "fino" {
		return g.stepEnd(in)
	}
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
		o.Events.Create = func(r *gm.Room, self *gm.Instance) { preparu(r) } // room start
		o.Events.Key = map[gm.Key]gm.Handler{gm.KeyPgDn: func(r *gm.Room, self *gm.Instance) {
			r.G.Score = -10000
			r.Goto = "@next"
		}}
	}
	// regionilo (the *_bildo rooms): room start sets alarm 0 = 1; alarm 0
	// turns the background's scroll (right → down → left → up → right…) and
	// comes back after 32, 64 or 96 steps.
	objectEvents["regionilo"] = func(a *Assets, o *gm.Object) {
		o.Events.Create = func(r *gm.Room, self *gm.Instance) { // room start
			self.Alarm[0] = 1
			preparu(r)
		}
		o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) {
			v := r.G.Vars
			h, vs := v["background_hspeed"], v["background_vspeed"]
			switch h {
			case 0:
				if vs == 1 {
					h, vs = -1, 0
				} else {
					h, vs = 1, 0
				}
			case 1:
				h, vs = 0, 1
			case -1:
				h, vs = 0, -1
			}
			v["background_hspeed"], v["background_vspeed"] = h, vs
			self.Alarm[0] = int(r.Random(3)+1) * 32
		}
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

// stepEnd is the end room: typing the name until Enter (the dialog stops
// the game), then any key or click → komenco.
func (g *Game) stepEnd(in gm.Input) error {
	if g.asking {
		if t, ok := in.(Typer); ok {
			g.G.Type(t.Typed())
		}
		if in.Pressed(gm.KeyEnter) {
			g.Highscores.Add(g.G.KeyboardString, g.G.Score)
			if g.SaveHighscores != nil {
				g.SaveHighscores(g.Highscores)
			}
			g.asking = false
		}
		return nil
	}
	g.World.Room.Step(in)
	if a, ok := in.(AnyKey); ok && a.AnyPressed() {
		if g.OnTitle != nil {
			g.OnTitle()
			return nil
		}
		return g.enter(0)
	}
	return nil
}

// Draw renders the room, and the end room's table.
func (g *Game) Draw(dst *image.RGBA) {
	g.World.Draw(dst)
	if RoomOrder[g.Room] == "fino" {
		for _, f := range g.World.Room.All("finigxu") {
			g.World.drawEnd(dst, int(f.X), int(f.Y), g.Highscores, g.asking, g.G.KeyboardString)
		}
	}
}

// preparu's music (regulilo's and regionilo's room start): unless muziko
// already plays (or MUTE, false as released), the CD's tracks when an audio CD
// is in, else muziko looped. A remake has no CD, so: muziko, looped, carrying
// on from room to room. (The background colour part is regionColour.)
func preparu(r *gm.Room) {
	if !r.G.Playing("muziko") {
		r.G.Play("@music:muziko")
	}
}
