# Pacmanio — loop (PRIVATE: contains stock assets; no remote, never publish without Joop's say-so)

One bounded chunk (~10 min) per tick, top unchecked item first. The recovered 2006 project is in `pacmanio/`
(20 rooms, 17 objects as GML, sprites, sounds) and the older preview in `pmv/`. Joop: Pacmanio "has been more
beautiful from the start and includes some of my own pixel art of the time. It still uses stock assets as well."

## Non-negotiables

- The recovered code and rooms are the truth; nothing is guessed. Where the code is unclear, note it as a
  question for Joop (no Crostini needed).
- No game code for a remake until Joop has approved a design (item 4). Commit locally each tick; no remote.
- Keep the recovered files as extracted (they are Joop's lost source); write new files beside them.

## Queue

1. [x] docs/PACMANIO-RULES.md from pacmanio/objects/*.gml and scripts: every object and event in step order
   (movement, ghosts/"afraid" mode, coins, bonus items, lives, scoring, level sets ar_/ls_/pu_, bildo rooms, fino).
2. [ ] Asset inventory: which sprites/sounds/backgrounds are Joop's own pixel art and which look like stock
   (GM 5 bundled resources) — a table for Joop to confirm; nothing removed.
3. [ ] Compare pmv/ (preview) with pacmanio/: what changed between them.
4. [ ] docs/DESIGN.md: a remake proposal on the Stelo model (Go + Ebitengine, original levels 1:1, rules and
   level-format docs, editor, packs), with what stock assets would need replacing → Joop. STOP there until yes.
5. [ ] (after the yes) The wasm build goes to ~/Projects/pacmanio-web, served PRIVATELY at
   https://kiefte.eu/pacmanio/ behind the goloco login (Joop 2026-11-01: 'add wasm pacmanio alongside goloco').
   Joop runs `sudo bash ~/Projects/pacmanio-recovered/scripts/caddy-pacmanio.sh` once (a placeholder page is
   there now); later deploys just replace the folder (atomic swap, chmod -R o+rX).

## Questions for Joop

- (none yet; the open engine points are in docs/PACMANIO-RULES.md §9)

## Blocked

## Tick log
- 2026-11-01 tick 1 (item 1) — docs/PACMANIO-RULES.md from all 17 objects + 3 scripts + rooms: Pacman, ghosts (adapt_direction; murox = Pacman-only wall, muroy = ghost-only), pellets (160+40 steps), fruit cycle (200 of every 700 steps), wall autotiling table, HUD corners, regions (pu/ar/ls), title-as-level, highscores, CD-or-muziko music. FOUND for item 2: the help text says 'surbaze de kodo de Mark Overmars; kelkaj bildoj estas de Billy McQuown' — GM's Pac-Man example is the stock base. DEBUG=true in the release: PageDown skips a level for −10000. Open: scared-ghost parent collision [engine].
- 2026-11-01 (between ticks) — scripts/caddy-pacmanio.sh written (reuses the goloco block's basic_auth line at runtime; nothing secret in the repo), placeholder ~/Projects/pacmanio-web/index.html; Joop to run the sudo step.
