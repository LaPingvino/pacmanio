// object regionilo (id 14)
// sprite=none solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Alarm 0
// --- code action ---
switch (background_hspeed) {
  case 0:
  switch (background_vspeed) {
    case 1:
        background_hspeed=-1;
        background_vspeed=0;
        break;
    default:
        background_hspeed=1;
        background_vspeed=0;
        break;
  }
  break;
  case 1:
    background_hspeed=0;
    background_vspeed=1;
    break;
  case -1:
    background_hspeed=0;
    background_vspeed=-1;
    break;
}

// --- end code ---
action_set_alarm(floor(random(3)+1)*32, 0);

#event Other: room end
// --- code action ---
global.chambro=cd_track();
// --- end code ---
if action_1_813()
action_1_809();

#event Other: room start
action_execute_script(preparu, 0, 0, 0, 0, 0);
action_set_alarm(1, 0);

#event Draw
action_draw_font("Trebuchet MS",36,16777215,1,1,0,0);
action_draw_text(''+ global.regiono, 0, 0)  /*relative*/;
action_draw_font("Trebuchet MS",36,0,1,1,0,0);
action_draw_text(''+ global.regiono, 1, 1)  /*relative*/;
