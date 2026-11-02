# Pacmanio

A Pac-Man game in Esperanto that Joop Kiefte made in 2006 with Game Maker 5.3, recovered from its own
executable and remade in Go with [Ebitengine](https://ebitengine.org/): the same rooms, rules and timing as the
original, playable on the desktop and in the browser.

**Play it in the browser:** https://lapingvino.github.io/pacmanio/ (arrow keys; on a phone, the on-screen pads).

## What is here

- `pacmanio/` — the 2006 project as recovered from `pacmanio.exe` (rooms, objects as pseudo-GML, sprites,
  backgrounds, sounds, scripts), with [gm53-reader](https://github.com/LaPingvino/gm53-reader).
  `pmv/` — an earlier preview build, recovered the same way.
- `docs/PACMANIO-RULES.md` — how the game works, read from the recovered code; the remake is written from it.
- `cmd/pacmanio`, `internal/` — the remake (Go). `levels/` — the rooms as plain-text level files.

Build and run: `go run ./cmd/pacmanio` (desktop), or for the web
`GOOS=js GOARCH=wasm go build -o pacmanio.wasm ./cmd/pacmanio` with `web/index.html` and Go's `wasm_exec.js`.

## Licence and the 2006 assets — please read

- **The code** (everything except the folders below) is MIT-licensed: see `LICENSE`.
- **The game's art, sounds and music** (`pacmanio/`, `pmv/`, and what is derived from them) are **not** covered
  by that licence. Pacmanio grew out of the Pac-Man example that came with Game Maker (by Mark Overmars, with
  images by Billy McQuown), made for Game Maker users to build their own games on; the game's own help text
  credits both. Next to that stock material are Joop's own pixel art (the houses, roads and hedges used as
  walls, the title and highscore pictures) and two MIDI tunes whose origin is not recorded.
  `docs/ASSETS.md` lists which is which as far as known. These are published here to preserve a 2006 hobby
  game; if you hold rights to any of it and want it removed or credited differently, please open an issue and
  it will be taken care of.
- **Fonts**: Liberation (SIL Open Font License), in `fonts/` with their licence.
