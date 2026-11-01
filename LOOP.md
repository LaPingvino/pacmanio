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
7. [ ] Pacman and the dots; ghosts; pellets; fruit; lives and HUD; regions, title and end rooms.
8. [ ] Music and sounds.
9. [ ] Mobile (as Stelo item 12).
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
