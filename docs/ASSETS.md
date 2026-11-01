# Pacmanio assets: Joop's own or stock?

For Joop to confirm. Joop's decision (2026-11-01): keep everything for now (private build); later replace stock
items with free-licensed or newly made ones, and keep Joop's own pixel art. The game's help text says it is
"surbaze de kodo de Mark Overmars" and that "kelkaj bildoj estas de Billy McQuown": Game Maker's Pac-Man example
(by Mark Overmars, with images by Billy McQuown) is the stock base. The **guess** column is mine, from what the
images show and how they're used. It is not checked against the GM example, which isn't on this machine.

The sprites are identical in `pacmanio/` and `pmv/`, so the preview doesn't show which came later.

## Sprites

| Sprite (frames) | What it is | Guess | Why |
|---|---|---|---|
| `pac_left/right/up/down` (4 each) | Pacman chomping | stock | classic Pac-Man frames, as in the GM example |
| `pac_stand` (4) | Pacman's smiling face (stopped; life icon) | unsure | not in classic Pac-Man; could be Joop's or the example's |
| `normal` (6) | green ghost with a red mouth | stock | ghost set with `afraid`/`inter` matches the example's three states |
| `afraid` (4) | light-blue frightened ghost | stock | as above |
| `inter` (4) | flashing ghost | stock | as above |
| `point_sprite` (1) | dot | stock | simple dot on dark ground, example style |
| `coin_sprite` (1) | power pellet: a yellow star | unsure | a star, not the classic pellet |
| `bonus_sprite` (4) | cherries, strawberry, banana, green apple | stock | clip-art style fruit, unlike the rest |
| `vivo_sprite` (1) | diamond (extra life) | unsure | |
| `muro_1` (16) | green tiles with grey channels (Akva Regiono walls) | unsure | could be the example's wall set |
| `muro_2` (16) | **houses and city blocks** (Urba Regiono walls) | **Joop's** | the houses Joop named as his own |
| `muro_3` (16) | grey roads, roundabouts, blue/yellow blocks (Libera Spaco walls) | Joop's (probable) | same hand as `muro_2`, made for the regions |
| `muro_x3` (16) | green hedges (`murox`, the walls only Pacman can't cross) | Joop's (probable) | as above |
| `muro_konstruo`, `murox_konstruo`, `muroy_konstruo` | editor placeholders (grey X, green X, red square) | Joop's | made for the room editor; not seen in play |

## Backgrounds

| Background | What it is | Guess |
|---|---|---|
| `komencigxo` (640×480) | the title screen: "Pacmanio — Projekto de iXio Esperantaj Ludoj", brick wall, big Pacman, four rings | Joop's (the brick texture and lettering effects may come from a graphics program) |
| `regiona` (64×64) | a frightened ghost on olive, tiled behind the region titles | stock-derived (the `afraid` ghost) |
| `poentegistoj` (256×64) | "Pacmanio Poentegistoj" header for the highscores | Joop's |

## Sounds

| Sound | Kind | Guess |
|---|---|---|
| `beeb`, `dead`, `bonus`, `won`, `click`, `catch` | short WAV effects | stock (the example's sound set) |
| `muziko` | MIDI, 19.6 KB (the level music when no CD) | unsure (not in `pmv`) |
| `finsono` | MIDI, 539 bytes (title and end screens) | unsure |

## Fonts

Drawn with Windows fonts (MS Sans Serif, Trebuchet MS, Courier New): the remake needs free replacements in any
case.

## Questions for Joop

1. Which of the "unsure" ones are yours: `pac_stand`, the star pellet, the diamond, `muro_1`, the two MIDI tunes?
2. Is `muro_1` (the green channel tiles) yours or the example's?
