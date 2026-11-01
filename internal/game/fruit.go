package game

import "pacmanio/internal/gm"

// The fruit (pacmanio/objects/bonus.gml) and the extra life (vivo.gml).
// Game Maker 5: a question action (here "is there no regulilo?", negated)
// governs only the action right after it, so the alarms are set in every
// room; only the moves depend on being in a level.

func init() {
	objectEvents["bonus"] = fruitEvents
	objectEvents["vivo"] = func(a *Assets, o *gm.Object) {
		o.Events.Collision = map[string]gm.CollisionHandler{
			"pacman": func(r *gm.Room, self, other *gm.Instance) {
				r.G.Play("won")
				r.G.AddLives(1)
				self.Destroyed = true
			},
		}
	}
}

func inLevel(r *gm.Room) bool { return len(r.All("regulilo")) != 0 }

// hide is action_move_to(-1000, 0) when in a level.
func hideFruit(r *gm.Room, self *gm.Instance) {
	if inLevel(r) {
		self.X, self.Y = -1000, 0
	}
}

// picture is image_single = floor(random(4)).
func picture(r *gm.Room, self *gm.Instance) {
	self.ImageIndex, self.ImageSpeed = float64(int(r.Random(4))), 0
}

func fruitEvents(a *Assets, o *gm.Object) {
	o.Events.Create = func(r *gm.Room, self *gm.Instance) {
		picture(r, self)
		hideFruit(r, self)
		self.Alarm[0] = 500
	}
	o.Events.Alarm[1] = func(r *gm.Room, self *gm.Instance) {
		hideFruit(r, self)
		self.Alarm[0] = 500
	}
	o.Events.Alarm[0] = func(r *gm.Room, self *gm.Instance) {
		picture(r, self)
		if inLevel(r) { // action_move_start
			self.X, self.Y = self.XStart, self.YStart
		}
		self.Alarm[1] = 200
	}
	o.Events.Collision = map[string]gm.CollisionHandler{
		// eaten: hidden whether in a level or not (no question before it)
		"pacman": func(r *gm.Room, self, other *gm.Instance) {
			r.G.Play("bonus")
			r.G.Score += 200
			self.X, self.Y = -1000, 0
			self.Alarm[0] = 500
		},
	}
}
