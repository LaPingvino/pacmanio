// This script adapts the direction of the monster
// when it comes at a possible crossing
{
  with (murox) {solid=false;}
  with (muroy) {solid=true;}
  if (hspeed == 0)
  {
    if (random(3)<1 && place_free(x-4,y))
      { hspeed = -4; vspeed = 0;}
    if (random(3)<1 && place_free(x+4,y))
      { hspeed =  4; vspeed = 0;}
  }
  else
  {
    if (random(3)<1 && place_free(x,y-4))
      { hspeed = 0; vspeed = -4;}
    if (random(3)<1 && place_free(x,y+4))
      { hspeed = 0; vspeed =  4;}
  }
  with (murox) {solid=true;}
  with (muroy) {solid=false;}
}
