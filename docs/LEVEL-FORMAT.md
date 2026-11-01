# Pacmanio level files

`levels/<room>.txt`, converted from `pacmanio/project.json` by `go run ./cmd/convert`; `go test ./...` checks
that all 20 rooms read back with the same settings and instances and that `levels/` is current.

The format is the Stelo remake's (its docs/LEVEL-FORMAT.md: header lines, `grid`, `instance`/`text` lines,
`end`) with `pacmanio-level 1` as the first line and **no ladder layer**. Grid characters:

| | Object | | Object |
|---|---|---|---|
| `#` | wall placeholder (`konstruu_muron`, becomes `muro` when created) | `o` | dot (`point`) |
| `W` | wall placed directly (`muro`) | `O` | power pellet (`pil`) |
| `x` | Pacman-only wall placeholder (`konstruu_murox` → `murox`) | `M` | ghost (`monster`) |
| `y` | ghost-only wall (`muroy`, invisible) | `B` | fruit (`bonus`) |
| `@` | Pacman | `V` | extra life (`vivo`) |

Control objects (`regulilo`, `regionilo`, `komencu`, `finigxu`) and a second object in an occupied cell are
`instance OBJECT X Y` lines.
