// Command render draws a room to a PNG (for checks).
package main

import (
	"flag"
	"image"
	"image/png"
	"log"
	"os"

	"pacmanio"
	"pacmanio/internal/game"
)

func main() {
	name := flag.String("room", "pu_1a_nivelo", "room")
	out := flag.String("out", "room.png", "output PNG")
	flag.Parse()
	a, err := game.LoadAssets(pacmanio.Files)
	if err != nil {
		log.Fatal(err)
	}
	r, err := game.LoadWorld(pacmanio.Files, a, *name)
	if err != nil {
		log.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, game.ScreenW, game.ScreenH))
	r.Draw(img)
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
