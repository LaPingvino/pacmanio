// object muro (id 3)
// sprite=muro_1 solid=True visible=True depth=0 persistent=False parent=none mask=none

#event Other: room start
// --- code action ---
switch (global.murigxu) {
case 1:sprite_index=muro_1;break;
case 2:sprite_index=muro_2;break;
case 3:sprite_index=muro_3;break;
}
// --- end code ---
// --- code action ---
if (place_free(x-32,y)) {
    if (place_free(x,y-32)) {
        if (place_free(x+32,y)) {
            if (place_free(x,y+32)) {
image_single=0;    //0000
            } else {
image_single=2;    //0001
            }    
        } else {
            if (place_free(x,y+32)) {
image_single=4;    //0010
            } else {
image_single=13;   //0011
            }   
        }
    } else {
        if (place_free(x+32,y)) {
            if (place_free(x,y+32)) {
image_single=1;    //0100
            } else {
image_single=14;    //0101
            }       
        } else {
            if (place_free(x,y+32)) {
image_single=10;    //0110
            } else {
image_single=9;    //0111
            }   
        }
    }
} else {
    if (place_free(x,y-32)) {
        if (place_free(x+32,y)) {
            if (place_free(x,y+32)) {
image_single=3;    //1000
            } else {
image_single=12;    //1001
            }       
        } else {
            if (place_free(x,y+32)) {
image_single=15;    //1010
            } else {
image_single=8;    //1011
            }   
        }    
    } else {
        if (place_free(x+32,y)) {
            if (place_free(x,y+32)) {
image_single=11;    //1100
            } else {
image_single=7;    //1101
            }       
        } else {
            if (place_free(x,y+32)) {
image_single=6;    //1110
            } else {
image_single=5;    //1111
            }   
        }
    }
}
// --- end code ---
