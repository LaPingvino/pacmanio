// object vivo (id 13)
// sprite=vivo_sprite solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Collision with pacman
action_sound(won, 0);
action_set_life(1)  /*relative*/;
action_kill_object();
