// object regulilo (id 8)
// sprite=none solid=False visible=True depth=-10 persistent=False parent=none mask=none

#event Keyboard <34>
// --- code action ---
if (DEBUG) {
score=-10000;
room_goto_next();
}
// --- end code ---

#event Other: room end
// --- code action ---
global.chambro=cd_track();
// --- end code ---
if action_1_813()
action_1_809();

#event Other: room start
action_execute_script(preparu, 0, 0, 0, 0, 0);

#event Draw
action_draw_font("MS Sans Serif",9,0,1,0,0,0);
// --- code action ---
brush_color =background_color;
pen_color=background_color;
font_color=global.kontrauxkoloro;
// --- end code ---
action_draw_rectangle(view_left[0], view_top[0], view_left[0]+96, view_top[0]+64);
action_draw_text('P A C M A N I O#' + room_caption, view_left[0]+3, view_top[0]+3);
action_draw_rectangle(view_left[0]+640-96, view_top[0]+0, view_left[0]+640, view_top[0]+64);
action_draw_score(view_left[0]+640-93, view_top[0]+3, 'Poentoj:#');
action_draw_rectangle(view_left[0]+0, view_top[0]+480-64, view_left[0]+96, view_top[0]+480);
action_draw_text(Programisto:##Joop KIEFTE, view_left[0]+3, view_top[0]+480-61);
action_draw_rectangle(view_left[0]+640-96, view_top[0]+480-64, view_left[0]+640, view_top[0]+480);
action_draw_text(Restaj vivoj:, view_left[0]+640-93, view_top[0]+480-61);
if !action_if_life(3, 2)
action_draw_life_images(view_left[0]+640-96, view_top[0]+480-32, pac_stand);
else
{
    action_draw_sprite(pac_stand, view_left[0]+640-96, view_top[0]+480-32, -1);
    action_draw_font("Courier New",18,65535,1,0,0,0);
    action_draw_life(view_left[0]+640-64, view_top[0]+480-32, '');
}
