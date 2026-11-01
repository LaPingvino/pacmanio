// object finigxu (id 12)
// sprite=none solid=False visible=True depth=0 persistent=False parent=none mask=none

#event Create
// --- code action ---
if (highscore_value(10)<score) {
highscore_add(get_string("Vi estas poentegisto! Kio fakte estas via nomo?","Pacman"),score);
}
// --- end code ---
action_sound(finsono, 1);

#event Step (normal)
// --- code action ---
if (mouse_button==mb_any || keyboard_key!=0) {room_goto(komenco);}
// --- end code ---

#event Other: room end
action_end_sound(finsono);

#event Draw
action_draw_font("Courier New",14,65535,1,0,0,0);
// --- code action ---
draw_highscore(x,y,x+256,y+256);
// --- end code ---
action_draw_text(Premu klavon por ekludi..., 0, 288)  /*relative*/;
