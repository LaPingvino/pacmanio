if (cd_present()) {
    if (not variable_global_exists("chambro")) {global.chambro=1;} else {
    global.chambro+=1;}
    if (global.chambro>cd_number()) {global.chambro=1;}
}

switch (global.murigxu) {
    case 1:background_color=c_blue;break;
    case 2:background_color=c_ltgray;break;
    case 3:background_color=c_dkgray;break;
}

if (not sound_isplaying(muziko) && not MUTE) {
    cd_init();
    if (cd_present()) {
        if (not cd_playing()) {cd_play(global.chambro,1000);}
    } else {
        sound_loop(muziko);
    }
}