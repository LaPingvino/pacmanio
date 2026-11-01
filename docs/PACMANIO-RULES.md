# How Pacmanio works

Read from the recovered 2006 project in `pacmanio/` (objects/*.gml, scripts/*.gml, rooms/*.json,
project.json). Behaviour that comes from **Game Maker 5.3 itself** rather than from Pacmanio's code is marked
**[engine]**; it follows the same step order as Stelo (see the Stelo repo's docs/STELO-RULES.md §2): alarms → keys
→ step → movement → outside-room → collisions (put back first if the other is solid) → end step → draw.

The game's own help text (`game_information.rtf`): "Farita de Joop KIEFTE, surbaze de kodo de Mark Overmars;
kelkaj bildoj estas de Billy McQuown". The `outside_wrap` script is credited to Mark Overmars (2001) in its
header: Pacmanio grew out of Game Maker's Pac-Man example.

## 1. Time, space, constants

- Level rooms run at **20 steps/s**, the region title rooms (`*_bildo`) at **15**, `komenco` and `fino` at **30**.
- Rooms are 640×480, except `pu_4a_nivelo` (640×960), `ar_4a_nivelo` (1280×960) and `ls_4a_nivelo`
  (1280×480). No room has views enabled [so the larger rooms show only their top-left 640×480; to check in play].
- Grid: 32 px. All sprites are 32×32, origin top-left. Bounding boxes (x left–right, y top–bottom):
  Pacman 1/2–29/30 × 1/2–29/30 (per direction), ghosts 1–30 × 0–31, dot (`point_sprite`) 12–18 × 12–18,
  power pellet (`coin_sprite`) 2–29 × 0–31, walls/bonus 0–31 × 0–31, life (`vivo_sprite`) 1–29 × 3–29.
- Constants: `DEBUG = true`, `MUTE = false` (as released).

## 2. Pacman (`pacman`)

- Animation speed 0.5 (`image_speed`). Game start: lives 3, score 0.
- **Keys, while held**: for each arrow, if the square 4 px that way is free of solids, move that way at speed 4
  and show `pac_left/right/up/down`. (So a key held against a wall does nothing, and Pacman keeps going in the
  old direction until it can turn.)
- **Each step**: when no `point` is left: sound `won`, pause 2 s, next room.
- **Collisions**:

| With | What happens |
|---|---|
| `muro` (wall, solid) | [engine] put back; stop; sprite `pac_stand` |
| `murox` (wall, solid) | back to the previous position; stop; `pac_stand` |
| `point` (dot) | the dot disappears, sound `click`, score +10 |
| `pil` (power pellet) | sound `beeb`, the pellet disappears; every `scared` becomes `monster` and then every `monster` becomes `scared` (so the timer restarts); every `scared` gets alarm 0 = 160 steps |
| `scared` (frightened ghost) | sound `catch`, the ghost goes back to its start and becomes `monster`, score +100 |
| `monster` | sound `dead`, pause 1.5 s, every ghost (monster and scared) back to its start, Pacman stops, back to his start, lives −1 |
| `bonus` (fruit) | see §4 (handled by the fruit) |
| `vivo` (extra life) | handled by the life: sound `won`, lives +1, it disappears |

- **Outside the room**: `outside_wrap` — leaving on one side brings it in on the other (only in the direction it
  is moving).
- **No more lives** → room `fino`.

## 3. Ghosts

- **`monster`** (sprite `normal`, 6 frames): starts at speed 4 in a random one of up/down/left/right. Each step,
  when aligned to the 32 px grid, runs `adapt_direction`. On `muro` [engine: put back first] or `muroy`: reverse
  both directions, then `adapt_direction`. Wraps around the room edges like Pacman.
- **`adapt_direction`**: for this check `murox` is not solid and `muroy` is. If moving vertically (hspeed 0): with
  chance 1/3, if free 4 px left, go left at 4; then with chance 1/3, if free 4 px right, go right at 4. If moving
  horizontally: the same with up, then down. (A later success overrides an earlier one.) Then `murox` is solid
  again and `muroy` not.
  So **`murox`** blocks Pacman but not ghosts, and **`muroy`** (invisible) blocks ghosts but not Pacman.
- **`scared`** (child of `monster`, sprite `afraid`): after alarm 0 (160 steps) becomes `intermonster` with
  alarm 40.
- **`intermonster`** (child of `scared`, sprite `inter`, the flashing one): after 40 steps becomes `monster`.
- **[engine] open point**: `scared` is a child of `monster`, and Pacman has a collision event with each. Whether
  Game Maker 5.3 runs both events when Pacman touches a scared ghost (and in which order) decides whether eating
  a ghost can also kill Pacman. To be checked by playing the original.

## 4. Fruit (`bonus`) and extra lives

- A fruit picks a random picture (frame 0–3). In a level (when a `regulilo` is present) it hides at (−1000, 0) and
  sets alarm 0 = 500 steps. Alarm 0: new random picture, back to its place, alarm 1 = 200 steps. Alarm 1: hide
  again, alarm 0 = 500. So it shows for 200 steps every 700.
- Eaten: sound `bonus`, score +200, hidden, alarm 0 = 500.
- Outside levels (title and end rooms, no `regulilo`) fruits stay where they were placed.
- `vivo`: an extra life (see §2).

## 5. Walls

- `konstruu_muron` / `konstruu_murox` are placeholders used in the room editor; they turn into `muro` / `murox`
  when created.
- At room start every `muro` takes the region's tile set (`global.murigxu` 1/2/3 → `muro_1/2/3`), and `muro`
  and `murox` pick their frame from which neighbours (32 px left, up, right, down) hold a solid:

| left up right down solid? | frame | | left up right down | frame |
|---|---|---|---|---|
| none | 0 | | left | 3 |
| down | 2 | | left, down | 12 |
| right | 4 | | left, right | 15 |
| right, down | 13 | | left, right, down | 8 |
| up | 1 | | left, up | 11 |
| up, down | 14 | | left, up, down | 7 |
| up, right | 10 | | left, up, right | 6 |
| up, right, down | 9 | | all four | 5 |

## 6. The level screen (`regulilo`)

- Room start: `preparu` (music and background colour, §8).
- Draws four 96×64 boxes in the corners of the view, filled with the background colour, with text (MS Sans Serif
  9, bold) in the region's text colour (`global.kontrauxkoloro`):
  top-left "P A C M A N I O" and the room caption; top-right "Poentoj:" and the score; bottom-left "Programisto:"
  and "Joop KIEFTE"; bottom-right "Restaj vivoj:" and the lives: up to 3 as `pac_stand` icons from
  (view right − 96, view bottom − 32); more than 3 as one icon and the number (Courier New 18, bold, yellow).
- **Page Down** (key 34), because `DEBUG` is true: score −10000 and the next room.

## 7. Rooms

Order: `komenco` → (`pu_bildo`, pu_1a … pu_4a, `pu_fina_nivelo`) → (`ar_bildo`, ar_1a … ar_4a, `ar_fina_nivelo`)
→ (`ls_bildo`, ls_1a … ls_4a, `ls_fina_nivelo`) → `fino`.

| Region | Name shown | Walls | Background | Text colour |
|---|---|---|---|---|
| pu | Urba Regiono | `muro_2` | light grey | black |
| ar | Akva Regiono | `muro_1` | blue | yellow |
| ls | Libera Spaco | `muro_3` | dark grey | light grey |

- **`komenco`** (title, also a tiny level): `komencu` stops the music, sets score 0 and lives 3, loops `finsono`,
  and after 450 steps goes to `fino` (the highscores) with transition 14. A mouse click quits the game. Pacman is
  playable here: eating the room's one dot starts the game (next room). Sets `global.murigxu = 1`.
- **`*_bildo`** (region title): shows the region name (`regionilo`: Trebuchet MS 36, drawn white, then black 1 px
  down-right). The background scrolls, changing direction every 32, 64 or 96 steps (right → down → left → up …).
  One dot to eat to go on.
- **Levels**: eat every dot. Room captions name the level (e.g. "Jen fantomoj!" for pu_1a).
- **`fino`** ("Fino de l' ludo"): `finigxu` — if the score beats the 10th highscore, asks "Vi estas poentegisto!
  Kio fakte estas via nomo?" (default "Pacman") and adds it; loops `finsono`; draws the highscores (Courier New
  14, yellow) and "Premu klavon por ekludi..." 288 px below; any key or click → `komenco`.

## 8. Music

`preparu` (each level and region room start): if an audio CD is present, plays its tracks in turn
(`global.chambro`, next track each room); otherwise loops `muziko`. Room end remembers the CD track.

## 9. Open points

- Scared-ghost collision and parent events (§3).
- What the larger rooms show without views (§1).
- `pmv/` (the preview) is compared in queue item 3.
