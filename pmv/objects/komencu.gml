// object komencu (id 10)
// sprite=none solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
action_set_score(0);
action_set_life(3);
if !action_if(MUTE)
action_sound(finsono, 1);
action_set_alarm(450, 0);

#event Alarm 0
action_another_room(fino, 14);

#event Step (normal)
// --- code action ---
if (mouse_button==mb_any) {game_end();}
// --- end code ---

#event Other: room end
action_end_sound(finsono);
