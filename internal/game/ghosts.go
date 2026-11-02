package game

import "pacmanio/internal/gm"

// The ghosts (pacmanio/objects/monster.gml, scared.gml, intermonster.gml,
// scripts/adapt_direction.gml) and the power pellet's effect.
// Pacman touching a ghost is pacman's collision events (objects.go): a
// scared one is eaten, a monster costs a life (Joop 2026-11-02, settling
// docs/PACMANIO-RULES.md §3's engine point).

func init() {
	objectEvents["monster"] = monsterEvents
	objectEvents["scared"] = func(a *Assets, o *gm.Object) {
		// Alarm 0 (160 steps after a pellet): flashing, for 40 more steps.
		o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) {
			r.Change(self, a.Object("intermonster"), false)
			for _, i := range r.All("intermonster") {
				i.Alarm[0] = 40
			}
		}
	}
	objectEvents["intermonster"] = func(a *Assets, o *gm.Object) {
		o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) { r.Change(self, a.Object("monster"), false) }
	}
}

// monsterEvents is monster.gml (inherited by scared and intermonster).
func monsterEvents(a *Assets, o *gm.Object) {
	// Create: action_move('010101010', 4) — one of down, left, right, up.
	o.Events.Create = func(r *gm.Room, self *gm.Instance) {
		dirs := [4][2]float64{{0, 4}, {-4, 0}, {4, 0}, {0, -4}}
		d := dirs[int(r.Random(4))]
		self.HSpeed, self.VSpeed = d[0], d[1]
	}
	// Step: on the 32 px grid (action_if_aligned(32, 32)), adapt_direction.
	o.Events.Step = func(r *gm.Room, self *gm.Instance) {
		if int(self.X)%32 == 0 && int(self.Y)%32 == 0 && self.X == float64(int(self.X)) && self.Y == float64(int(self.Y)) {
			adaptDirection(r, self)
		}
	}
	turn := func(r *gm.Room, self, other *gm.Instance) {
		self.HSpeed, self.VSpeed = -self.HSpeed, -self.VSpeed
		adaptDirection(r, self)
	}
	o.Events.Collision = map[string]gm.CollisionHandler{"muro": turn, "muroy": turn}
	o.Events.OutsideRoom = outsideWrap
}

// adaptDirection is the adapt_direction script: for the check murox isn't
// solid and muroy is; moving vertically, with chance 1/3 each it turns left
// (then right) if that way is free; moving horizontally, up (then down).
func adaptDirection(r *gm.Room, self *gm.Instance) {
	a := r.Instances
	setSolid := func(name string, v bool) {
		for _, i := range a {
			if i.Obj.Name == name {
				i.Solid = v
			}
		}
	}
	setSolid("murox", false)
	setSolid("muroy", true)
	if self.HSpeed == 0 {
		if r.Random(3) < 1 && r.PlaceFree(self, self.X-4, self.Y) {
			self.HSpeed, self.VSpeed = -4, 0
		}
		if r.Random(3) < 1 && r.PlaceFree(self, self.X+4, self.Y) {
			self.HSpeed, self.VSpeed = 4, 0
		}
	} else {
		if r.Random(3) < 1 && r.PlaceFree(self, self.X, self.Y-4) {
			self.HSpeed, self.VSpeed = 0, -4
		}
		if r.Random(3) < 1 && r.PlaceFree(self, self.X, self.Y+4) {
			self.HSpeed, self.VSpeed = 0, 4
		}
	}
	setSolid("murox", true)
	setSolid("muroy", false)
}

// pellet is pacman's collision with pil: beeb, the pellet goes, every scared
// ghost becomes a monster and then every monster (with its children) scared,
// with 160 steps on alarm 0.
func pellet(a *Assets) gm.CollisionHandler {
	return func(r *gm.Room, self, other *gm.Instance) {
		r.G.Play("beeb")
		other.Destroyed = true
		for _, g := range r.All("scared") {
			r.Change(g, a.Object("monster"), false)
		}
		for _, g := range r.All("monster") {
			r.Change(g, a.Object("scared"), false)
		}
		for _, g := range r.All("scared") {
			g.Alarm[0] = 160
		}
	}
}
