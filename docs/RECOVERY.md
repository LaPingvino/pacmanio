# How the 2006 project was recovered

Extracted on 2026-11-01 from Joop's own Game Maker 5.3a builds with
[gm53-reader](https://github.com/LaPingvino/gm53-reader) (no third-party code ran).

- `pacmanio/` — from `pacmanio.exe` (md5 cb6075439f28eb61a9ef6e3866eef1b0): 20 rooms (komenco = title, three
  level sets ar_/ls_/pu_ of 4 levels + a final level + a picture room "bildo" each, fino = end), 17 objects as
  pseudo-GML, 108 sprite frames, 3 backgrounds, 8 sounds, 3 scripts. Parsed to 2,120,208 of 2,129,965 bytes; the
  ~10 KB tail (probably the resource tree) is not read yet.
- `pmv/` — from `pmv.exe` (md5 1d349b5562d70068043a845466f4f9ae), probably an earlier preview: 14 rooms, 14 objects.
- `rooms-contact.png` — every room rendered with its own sprites.

To re-run: `python -I gm53_reader.py EXE OUT && python -I gm53_export.py OUT`.
