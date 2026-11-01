// object monster (id 1)
// sprite=normal solid=False visible=True depth=-2 persistent=False parent=none mask=none

#event Create
action_move('010101010', 4);

#event Step (normal)
if action_if_aligned(32, 32)
action_execute_script(adapt_direction, 0, 0, 0, 0, 0);

#event Collision with muroy
action_reverse_xdir();
action_reverse_ydir();
action_execute_script(adapt_direction, 0, 0, 0, 0, 0);

#event Collision with muro
action_reverse_xdir();
action_reverse_ydir();
action_execute_script(adapt_direction, 0, 0, 0, 0, 0);

#event Other: outside room
action_execute_script(outside_wrap, 0, 0, 0, 0, 0);
