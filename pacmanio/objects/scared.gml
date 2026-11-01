// object scared (id 2)
// sprite=afraid solid=False visible=True depth=-2 persistent=False parent=monster mask=none

#event Alarm 0
action_change_object(intermonster, 0);
with (intermonster) action_set_alarm(40, 0);
