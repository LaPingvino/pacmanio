// object pacman (id 0)
// sprite=pac_stand solid=False visible=True depth=-1 persistent=False parent=none mask=none

#event Create
image_speed = 0.5;

#event Step (normal)
if action_if_number(point, 0, 0)
{
    action_sound(won, 0);
    action_sleep(2000, 1);
    action_next_room(0);
}

#event Collision with pil
action_sound(beeb, 0);
with (other) action_kill_object();
with (scared) action_change_object(monster, 0);
with (monster) action_change_object(scared, 0);
with (scared) action_set_alarm(160, 0);

#event Collision with point
with (other) action_kill_object();
action_sound(click, 0);
action_set_score(10)  /*relative*/;

#event Collision with muro
action_move('000010000', 0);
action_set_sprite(pac_stand, 0);

#event Collision with scared
action_sound(catch, 0);
with (other) action_move_start();
with (other) action_change_object(monster, 0);
action_set_score(100)  /*relative*/;

#event Collision with monster
action_sound(dead, 0);
action_sleep(1500, 1);
with (monster) action_move_start();
with (scared) action_move_start();
action_move('000010000', 0);
action_set_sprite(pac_stand, 0);
action_move_start();
action_set_life(-1)  /*relative*/;

#event Keyboard <down>
if action_if_empty(0, 4, 0)  /*relative*/
{
    action_move('010000000', 4);
    action_set_sprite(pac_down, 0);
}

#event Keyboard <right>
if action_if_empty(4, 0, 0)  /*relative*/
{
    action_move('000001000', 4);
    action_set_sprite(pac_right, 0);
}

#event Keyboard <up>
if action_if_empty(0, -4, 0)  /*relative*/
{
    action_move('000000010', 4);
    action_set_sprite(pac_up, 0);
}

#event Keyboard <left>
if action_if_empty(-4, 0, 0)  /*relative*/
{
    action_move('000100000', 4);
    action_set_sprite(pac_left, 0);
}

#event Other: no more lives
action_another_room(fino, 0);

#event Other: game start
action_set_life(3);
action_set_score(0);

#event Other: outside room
action_execute_script(outside_wrap, 0, 0, 0, 0, 0);
