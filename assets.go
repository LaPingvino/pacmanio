// Package pacmanio holds the game's data: the recovered project's sprites,
// backgrounds and description, and the level files. PRIVATE (stock assets).
package pacmanio

import "embed"

// Files holds pacmanio/project.json, pacmanio/sprites/*.png,
// pacmanio/backgrounds/*.png and levels/*.txt.
//
//go:embed pacmanio/project.json pacmanio/sprites/*.png pacmanio/backgrounds/*.png levels/*.txt fonts/*.ttf
var Files embed.FS
