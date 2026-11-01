# The two Pacmanio builds

Both were saved on 15 April 2006 (Game Maker date 38823): `pacmanio.exe` at about 08:24, `pmv.exe` at about
22:05 the same day. `pmv` is not an older preview: it is a cut-down **windowed "work" edition** made from the
full game that evening. Its own description says so:

- `pmv` (settings information): "Ĉi tiu versio de la ludo ne estas tiom kompleta kiom la 'granda' versio. Ĉi
  tiu versio estas por ludado ekz. en laborejo..." — not as complete as the big version; for playing e.g. at
  work.
- `pmv` (help text): "ATENTU: Ĉi tiu fenestra versio de la ludo ne enhavas ĉiujn nivelojn, nur la unuajn
  10..." — this windowed version doesn't contain all the levels, only the first 10.
- `pacmanio` (settings information): "Atentu, ĉi tiu versio estas ankoraŭ ne kompleta. Atendu ĝis versio 1.0
  por kompleta ludo!" — not complete yet; wait for version 1.0.

## What `pmv` leaves out

| | `pacmanio` (the "big" version) | `pmv` (the work version) |
|---|---|---|
| Rooms | 20: title, Urba Regiono, Akva Regiono, **Libera Spaco** (6 rooms), end | 14: no Libera Spaco |
| Level music | `muziko.mid`, or the tracks of an audio CD if one is in the drive (`preparu`) | none: no `muziko`, no CD code |
| `murox` / `muroy` walls | yes (Pacman-only and ghost-only walls) and their collision events | absent (not needed without Libera Spaco) |
| Window | fullscreen, sets the screen resolution, custom loading image | windowed |
| `DEBUG` | true (Page Down skips a level for −10000) | false |
| Game settings version | 15 | 1 |

Everything else is the same: the 14 shared rooms are identical instance for instance, the sprites are
identical, and both have the same six sound effects and `finsono`.

So: **`pmv` lacks the levels and the music; `pacmanio` has both.** Joop remembered that one build lacked sound
and the other levels; by the files, `pacmanio` is complete on both counts. Its music depends on MIDI playback
(or an audio CD), which may simply not have played on the machine it was tried on.

## For the remake

`pacmanio` is the source: all 20 rooms, the music, the special walls. Nothing in `pmv` is missing from it.
