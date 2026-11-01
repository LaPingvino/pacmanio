package game

import "pacmanio/internal/gm"

// The objects' events, ported from pacmanio/objects/*.gml and scripts/*.gml.
var objectEvents = map[string]func(a *Assets, o *gm.Object){}

func init() {
	// The room-editor placeholders become walls when created
	// (action_change_object(muro, 1)).
	objectEvents["konstruu_muron"] = func(a *Assets, o *gm.Object) {
		o.Events.Create = func(r *gm.Room, self *gm.Instance) { r.Change(self, a.Object("muro"), true) }
	}
	objectEvents["konstruu_murox"] = func(a *Assets, o *gm.Object) {
		o.Events.Create = func(r *gm.Room, self *gm.Instance) { r.Change(self, a.Object("murox"), true) }
	}
	objectEvents["pacman"] = pacmanEvents
}

// stop is action_move('000010000', 0) and the standing sprite (frame kept,
// action_set_sprite(pac_stand, 0)).
func stop(self *gm.Instance) {
	self.HSpeed, self.VSpeed = 0, 0
	self.Sprite = "pac_stand"
}

// pacmanEvents is pacman.gml (dots, pellets, walls, keys, wrap). PARTIAL:
// touching ghosts (open engine point), fruit, extra life, lives (item 7).
func pacmanEvents(a *Assets, o *gm.Object) {
	o.Events.Create = func(r *gm.Room, self *gm.Instance) { self.ImageSpeed = 0.5 }
	// Step: no dots left → sound won, 2 s, next room.
	o.Events.Step = func(r *gm.Room, self *gm.Instance) {
		if len(r.All("point")) == 0 && r.Goto == "" {
			r.G.Play("won")
			r.SleepMS(2000, roomSpeed(r))
			r.Goto = "@next"
		}
	}
	// Keys held: if the square 4 px that way is free, move that way at 4.
	key := func(dx, dy float64, sprite string) gm.Handler {
		return func(r *gm.Room, self *gm.Instance) {
			if r.PlaceFree(self, self.X+dx, self.Y+dy) {
				self.HSpeed, self.VSpeed = dx, dy
				self.Sprite = sprite // action_set_sprite(…, 0): frame kept
			}
		}
	}
	o.Events.Key = map[gm.Key]gm.Handler{
		gm.KeyDown: key(0, 4, "pac_down"), gm.KeyRight: key(4, 0, "pac_right"),
		gm.KeyUp: key(0, -4, "pac_up"), gm.KeyLeft: key(-4, 0, "pac_left"),
	}
	o.Events.Collision = map[string]gm.CollisionHandler{
		// muro is solid: [engine] put back first; then stop.
		"muro": func(r *gm.Room, self, other *gm.Instance) { stop(self) },
		// murox: back to the previous position, stop.
		"murox": func(r *gm.Room, self, other *gm.Instance) {
			self.X, self.Y = self.XPrevious, self.YPrevious
			stop(self)
		},
		"pil": pellet(a),
		"point": func(r *gm.Room, self, other *gm.Instance) {
			other.Destroyed = true
			r.G.Play("click")
			r.G.Score += 10
		},
	}
	o.Events.OutsideRoom = outsideWrap
}

// outsideWrap is the outside_wrap script (Mark Overmars, 2001): leaving the
// room on a side moving outwards brings the instance in on the other side.
// Every Pacmanio sprite has its origin at 0, 0 (sprite_xoffset = 0).
func outsideWrap(r *gm.Room, self *gm.Instance) {
	const size = 32 // sprite_width / sprite_height
	switch {
	case self.X < 0 && self.HSpeed < 0:
		self.X = float64(r.W)
	case self.X > float64(r.W) && self.HSpeed > 0:
		self.X = -size
	}
	switch {
	case self.Y < 0 && self.VSpeed < 0:
		self.Y = float64(r.H)
	case self.Y > float64(r.H) && self.VSpeed > 0:
		self.Y = -size
	}
}

// roomSpeed is the room's steps per second (20 in levels, 15 in the region
// rooms, 30 in komenco and fino).
func roomSpeed(r *gm.Room) int {
	if s, ok := r.G.Vars["room_speed"]; ok {
		return int(s)
	}
	return 20
}
