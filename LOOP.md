# Pacmanio — loop (PRIVATE: contains stock assets; no remote, never publish without Joop's say-so)

One bounded chunk (~10 min) per tick, top unchecked item first. The recovered 2006 project is in `pacmanio/`
(20 rooms, 17 objects as GML, sprites, sounds) and the older preview in `pmv/`. Joop: Pacmanio "has been more
beautiful from the start and includes some of my own pixel art of the time. It still uses stock assets as well."

## Non-negotiables

- The recovered code and rooms are the truth; nothing is guessed. Where the code is unclear, note it as a
  question for Joop (no Crostini needed).
- Design approved 2026-11-01 (docs/DESIGN.md). Commit locally each tick; no remote; the web build only at the private kiefte.eu/pacmanio/.
- Keep the recovered files as extracted (they are Joop's lost source); write new files beside them.

## Queue

1. [x] docs/PACMANIO-RULES.md from pacmanio/objects/*.gml and scripts: every object and event in step order
   (movement, ghosts/"afraid" mode, coins, bonus items, lives, scoring, level sets ar_/ls_/pu_, bildo rooms, fino).
2. [x] Asset inventory: which sprites/sounds/backgrounds are Joop's own pixel art and which look like stock
   (GM 5 bundled resources) — a table for Joop to confirm; nothing removed.
3. [x] Compare pmv/ (preview) with pacmanio/: what changed between them. Joop (2026-11-01): "I think the
   Pacmanio builds both lack some. One lacks sound, the other lacks levels." First look: pmv lacks the whole
   Libera Spaco region (ls_*: 6 rooms) and the `muziko` MIDI; the 14 rooms they share are identical, and both
   have the same six WAVs. So far pacmanio.exe looks like the superset — check objects/scripts/settings too,
   and whether pacmanio's sound fails at run time (CD-audio path in `preparu`, MIDI) rather than being absent.
4. [x] docs/DESIGN.md: a remake proposal on the Stelo model (Go + Ebitengine, original levels 1:1, rules and
   level-format docs, editor, packs), with what stock assets would need replacing → Joop. STOP there until yes.
5. [x] Level converter for the 20 rooms (round-trip test), reusing the Stelo format.
6. [x] Rooms drawn with the original sprites, wall autotiling (test against the room renders).
7. [~] Pacman and the dots (done); ghosts; pellets; fruit; lives and HUD; regions, title and end rooms.
8. [~] Music and sounds (effects done; the two MIDI tunes wait for soundfont-fluid, as Stelo's).
9. [~] Mobile (as Stelo item 12): done but for a real-phone try.
10. [ ] The wasm build goes to ~/Projects/pacmanio-web, served PRIVATELY at
   https://kiefte.eu/pacmanio/ behind the goloco login (Joop 2026-11-01: 'add wasm pacmanio alongside goloco').
   Joop runs `sudo bash ~/Projects/pacmanio-recovered/scripts/caddy-pacmanio.sh` once (a placeholder page is
   there now); later deploys just replace the folder (atomic swap, chmod -R o+rX).

## Decisions (Joop)

- 2026-11-01: DESIGN APPROVED ("I think there is nothing I can disagree with among what you propose"). On the
  three open points the original's behaviour is kept: DEBUG as released (PageDown skip), the big rooms without
  views as in the original, and the highscore name in an in-game text field (no Windows dialog exists).
- 2026-11-01: true Esperanto letters on screen instead of the x-system (as for Stelo).
- 2026-11-01: mobile support approved, as for Stelo.

- 2026-11-01 assets: "For now you can keep the same assets. If you can find a replacement for stock stuff, you
  can of course keep my own pixel art (houses used as walls etc...)" → the remake starts with the original
  sprites/sounds (private build only); stock items (GM Pac-Man example / Billy McQuown images) get replaced
  later by free-licensed or newly made ones, Joop's own pixel art stays. Item 2's inventory says which is which.

## Questions for Joop

- Do you remember: could Pacman eat the blue (scared) ghosts in Pacmanio, or did touching one still cost a life? (Settles the open engine point in PACMANIO-RULES §3: a scared ghost is also a 'monster', so both of Pacman's collision events see it.)


- docs/ASSETS.md §Questions: which 'unsure' assets are his (pac_stand, star pellet, diamond, muro_1, the two MIDIs).
- (engine open points: docs/PACMANIO-RULES.md §9)

## Blocked

## Tick log
- 2026-11-01 tick 1 (item 1) — docs/PACMANIO-RULES.md from all 17 objects + 3 scripts + rooms: Pacman, ghosts (adapt_direction; murox = Pacman-only wall, muroy = ghost-only), pellets (160+40 steps), fruit cycle (200 of every 700 steps), wall autotiling table, HUD corners, regions (pu/ar/ls), title-as-level, highscores, CD-or-muziko music. FOUND for item 2: the help text says 'surbaze de kodo de Mark Overmars; kelkaj bildoj estas de Billy McQuown' — GM's Pac-Man example is the stock base. DEBUG=true in the release: PageDown skips a level for −10000. Open: scared-ghost parent collision [engine].
- 2026-11-01 (between ticks) — scripts/caddy-pacmanio.sh written (reuses the goloco block's basic_auth line at runtime; nothing secret in the repo), placeholder ~/Projects/pacmanio-web/index.html; Joop to run the sudo step.
- 2026-11-01 (between ticks) — Joop's asset decision recorded under Decisions.
- 2026-11-01 tick 2 (item 2) — docs/ASSETS.md: every sprite/background/sound with a guess (own/stock/unsure) and why; Joop's: muro_2 houses (named by him), probably muro_3/muro_x3, title and highscore backgrounds, editor placeholders; stock: Pacman, ghosts, fruit, dot, the six WAVs; unsure listed as questions. Sprites identical in pmv; pmv was saved ~14 h AFTER pacmanio (same day, 2006-04-15) and lacks muziko → item 3.
- 2026-11-01 tick 3 (item 3) — docs/BUILDS.md: pmv is the cut-down windowed 'work' edition made the same evening (its own texts say so: 'por ludado ekz. en laborejo', 'nur la unuajn 10 nivelojn'): no Libera Spaco, no music (no muziko, no CD code), no murox/muroy, windowed, DEBUG false. pacmanio.exe is the superset → the remake's source. NEXT: item 4 (design for Joop).
- 2026-11-01 tick 4 (item 4) — docs/DESIGN.md written (Stelo model, shared GM engine core, private web build, assets per Joop's decision, three questions on deviations). Now WAITING on Joop's yes; next workable items only after it.
- 2026-11-01 — Joop approved the design; queue items 5-10 added.
- 2026-11-01 tick 5 (item 5) — Go module 'pacmanio'; internal/level adapted from Stelo's (pacmanio-level 1, no ladder layer, Pacmanio grid chars), cmd/convert → levels/*.txt (20 rooms), docs/LEVEL-FORMAT.md; test: all 20 round-trip, levels/ current. NEXT: item 6 (draw rooms, wall autotiling).
- 2026-11-01 tick 6 (item 6) — internal/game: assets (from Stelo's loader), rooms from levels/, placeholders → muro/murox, muro's 16-case autotile table (unit test, every case), region tile set and preparu's background colour (pu → muro_2/light grey, ar → muro_1/blue, ls → muro_3/dark grey; ASSUMPTION: fino keeps the last region), drawing by depth, top-left 640×480 (no views). cmd/render; pu/ar/ls level 1 rendered and checked. The room renders can't be compared pixel-wise: gm53_export drew the placeholders' own sprites. Fruit visible until the create events run (item 7).
- 2026-11-01 tick 7 (item 7) — rooms now run on a copy of Stelo's Game Maker core (internal/gm; Stelo's is internal to its module): objects with sprite box, solidity, depth, parents (resolved from object ids, which have gaps); placeholders change to muro/murox on create; muro's room start autotiles via place_free; pacman: held keys move at 4 when free (sprite kept frame), muro (solid put-back) / murox (back to previous) stop him, dots +10, outside_wrap, all dots eaten → 'won', 2 s, next room. Tests: 16 wall cases on the engine, Pacman eats dots going left and stops at the wall. NEXT: ghosts (adapt_direction), pellets, fruit, lives.
- 2026-11-01 tick 8 (item 7) — ghosts: random start among down/left/right/up at 4, adapt_direction on the 32 px grid (murox not solid, muroy solid during the check: gm instances now have their own solid variable), reverse + adapt on muro/muroy, wrap; pellets: scared (alarm 160) → intermonster (40) → monster; gm.Room.Random (fixed seed; GM's generator not reproduced). Test: a ghost roams 200 steps never inside a wall; the pellet cycle 160/40. NOT PORTED: Pacman touching ghosts (open engine point; asked Joop).
- 2026-11-01 tick 9 (item 7) — fruit (bonus.gml) and extra life (vivo.gml), ported literally: GM5's question action governs only the next action, so alarms run in every room and only the moves need a level (rules §4 corrected). Tests: hidden → at its place at 500 → hidden at 700; eaten +200; vivo +1 life. Render: komenco's fruits with random pictures. LEFT for 7: Pacman touching ghosts (waits on Joop's answer), lives/no more lives → fino, regulilo HUD.
- 2026-11-01 tick 10 (item 7) — regulilo's boxes: four corner boxes in the background colour, text in the region's colour (Liberation Sans Bold 9 for MS Sans Serif, Liberation Mono Bold 18 for Courier New; OFL, fonts/), caption, score, programmer, lives (icons up to 3, else icon + number). Found: regulilo is every level's last instance, so level fruits show for the first 699 steps (alarm 1 set inside alarm 0 counts down that same step; rules §4 fixed; test). Render: pu_1a and ar_1a. LEFT for 7: ghosts touching Pacman (waits on Joop), no more lives → fino, Page Down debug skip.
- 2026-11-01 tick 11 (item 7) — Game: rooms in project.json's order with shared score/lives, the region from komenco (1) and the *_bildo rooms (fino keeps it), room changes after the pause, no more lives → fino; regulilo's Page Down (DEBUG: score becomes −10000, next room); komencu (score 0, lives 3, finsono, alarm 450 → fino). Test: komenco → pu_bildo → pu_1a → Page Down → pu_2a → no lives → fino (region 2). Rules: Page Down sets the score (not subtracts). LEFT for 7: ghosts touching Pacman (waits on Joop); then item 8 (fino/highscores, region titles' scroll).
- 2026-11-01 tick 12 (item 7) — region title rooms: tiled room backgrounds with background_hspeed/vspeed (per room, start 0); regionilo: alarm 0 at room start (1 step) turns the scroll right → down → left → up, every 32/64/96 steps (random); its draw: the region name (global.regiono) in white, then black 1 px down-right, Liberation Sans Bold Italic 36 for Trebuchet MS 36 bold italic. Test: the scroll's turns; render: pu_bildo, ls_bildo. LEFT for 7: ghosts touching Pacman (waits on Joop), fino (finigxu: highscores, name entry).
- 2026-11-01 tick 13 (item 7) — the end room (finigxu): Game Maker's 10-entry table ("<nobody>" 0 assumed, as Stelo); a score above the 10th asks the name in an in-game field ("Pacman" filled in, Backspace edits — gm Type got Stelo's Backspace), in a background-coloured box centred on screen (the original's get_string was a Windows dialog); Enter adds it (SaveHighscores hook); the table in Liberation Mono Bold 14 yellow over (x, y)–(x+256, y+256) and "Premu klavon por ekludi..." at y+288; then any key/click → komenco. Test: 1230 → name 'Joop' first in the table → komenco. Render: question box and table. LEFT for 7: ghosts touching Pacman (waits on Joop); saving the table (store, with the desktop/web build, item 10).
- 2026-11-01 tick 14 (item 8) — the game is playable: cmd/pacmanio (Ebitengine, desktop + wasm; the room's own speed as TPS: 20/15/30; keyboard, Backspace, any key/click), internal/sound from Stelo (the six WAVs embedded, decoded at 44.1 kHz; test), web/index.html (minimal, as Stelo's without the phone keyboard yet). Headless wasm: the title room with the 2006 title art. PARTIAL: muziko/finsono MIDI (soundfont), highscore storage. NEXT: item 9 mobile or item 10 deploy (Joop's caddy script first).
- 2026-11-01 tick 15 (item 9) — mobile, from Stelo: arrow pads on touch screens (or ?pads), a tap = any key, the phone keyboard ('Klavaro', pacmanioTyped) for the highscore name with Backspace/Enter, no page zoom/scroll (web/index.html now Stelo's full page). Headless 844×390 ?pads: the title fills the height, pads in the corners, Klavaro top-right. LEFT: a real phone; the pads cover the bottom corners of the room (as in Stelo).
