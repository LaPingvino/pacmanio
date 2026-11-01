////////////////////////////////////////////////////////////////
// script name: outside_wrap
//     creator: Mark Overmars
//        date: May 26, 2001
//
// description: When an instance moves outside the room, this script
//              wraps it to the other side of the room.
//
//   arguments: None
//
//     remarks: Place in the Outside event of the object.
////////////////////////////////////////////////////////////////
{
  if (x < 0 && hspeed < 0) x = room_width + sprite_xoffset;
  if (x > room_width && hspeed > 0) x = -sprite_width + sprite_xoffset;
  if (y < 0 && vspeed < 0) y = room_height + sprite_yoffset;
  if (y > room_height && vspeed > 0) y = -sprite_height + sprite_yoffset;
}
