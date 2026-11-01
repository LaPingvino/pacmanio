# Pacmanio remake — design (proposal, waiting for Joop's yes)

On the Stelo model: the 2006 game first, faithful, then an editor on top. Shares its engine with the Stelo
remake.

## 1. Baseline: the original game

- Go + Ebitengine, desktop and WebAssembly. The web build is **private**: served at
  https://kiefte.eu/pacmanio/ behind the goloco login (`scripts/caddy-pacmanio.sh`), because the original
  assets include stock material.
- Plays `pacmanio.exe` (the full version, docs/BUILDS.md) as it was: the title room (eat the dot to start), the
  three regions (Urba Regiono, Akva Regiono, Libera Spaco) with their title rooms and 5 levels each, the end
  room with the highscores. All 20 rooms converted 1:1 from `pacmanio/rooms/*.json`.
- Rules exactly as docs/PACMANIO-RULES.md: 20 steps/s in levels (15 in region titles, 30 in title/end),
  Pacman's turning rule, `adapt_direction`, the pellet timers (160 + 40 steps), the fruit cycle, the wall
  autotiling, `murox`/`muroy`, wrap-around, the HUD corners.
- Its engine is the Stelo remake's Game Maker core (bounding boxes, place_free, motion, step order). If Joop
  chooses the GML interpreter for Stelo, Pacmanio runs its recovered GML on it too, with the same built-ins;
  otherwise the objects are ported to Go from the rules document.
- The open engine question (does eating a frightened ghost also trigger the parent `monster` collision?) is
  settled by the engine's parent-event rule, written once for both games.

## 2. Assets

- Start with all the original sprites, backgrounds and sounds (Joop, 2026-11-01).
- Later: replace the stock ones (docs/ASSETS.md: Pacman, ghosts, fruit, dot, the six WAVs) with free-licensed
  or newly made ones; keep Joop's own pixel art (the house walls and the other region tiles, title and
  highscore art).
- Music: `muziko.mid` and `finsono.mid` rendered to OGG for the web. The audio-CD path is dropped (there is no CD).
- Fonts: the Windows fonts (MS Sans Serif, Trebuchet MS, Courier New) replaced by free look-alikes.

## 3. Changes from the original that need Joop's word

- `DEBUG` is true in the release (Page Down skips a level for −10000 points): keep, or set false as in `pmv`?
- The rooms bigger than 640×480 have no views in the original, so part of them is off screen. Keep that, or
  add a view following Pacman?
- The highscore name prompt was a Windows dialog: an in-game text field instead.

## 4. Level files and editor (later)

The Stelo level format (a grid, extra layers, instance lines) with Pacmanio's characters; the same editor.

## 5. Order of work (after the yes)

1. Level converter for the 20 rooms (round-trip test).
2. Rooms drawn with the original sprites, wall autotiling (test: matches the room renders).
3. Pacman and the dots; then ghosts; pellets; fruit; lives and HUD; regions and title/end rooms.
4. Music and sounds; the private web build and deploy.
5. Asset replacement for the stock items; the editor.
