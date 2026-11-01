// object bonus (id 7)
// sprite=bonus_sprite solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
// --- code action ---
image_single=floor(random(4));
// --- end code ---
if !action_if_number(regulilo, 0, 0)
action_move_to(-1000, 0);
action_set_alarm(500, 0);

#event Alarm 1
if !action_if_number(regulilo, 0, 0)
action_move_to(-1000, 0);
action_set_alarm(500, 0);

#event Alarm 0
// --- code action ---
image_single=floor(random(4));
// --- end code ---
if !action_if_number(regulilo, 0, 0)
action_move_start();
action_set_alarm(200, 1);

#event Collision with pacman
action_sound(bonus, 0);
action_set_score(200)  /*relative*/;
action_move_to(-1000, 0);
action_set_alarm(500, 0);
