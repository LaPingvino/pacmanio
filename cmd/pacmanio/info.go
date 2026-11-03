package main

import (
	"image/color"

	"pacmanio/internal/game"
)

// The title's "Kiel ludi" (how to play) and "Pri" (about) pages — the texts of
// docs/TEXTS.md (UI step U3), in Esperanto (Pacmanio's one language). DRAFT
// texts, put in so Joop can correct them in place (his choice, 2026-11-04).
var infoPages = map[string][]string{
	"help": {
		"# Kiel ludi",
		"Manĝu ĉiujn punktojn de nivelo por iri al la sekva. Tri regionoj",
		"atendas vin — la Urba Regiono, la Akva Regiono kaj la Libera Spaco —,",
		"ĉiu kun kvar niveloj kaj fina nivelo.",
		"",
		"* Klavoj: sagoklavoj por iri · Esc aŭ P paŭzo · M sono · F11 plena ekrano",
		"* Telefono: la butonoj ĉe la randoj.",
		"Pacman turniĝas tuj kiam la vojo estas libera: tenu la klavon antaŭe.",
		"",
		"Punkto +10 · Potenca monero: la fantomoj timiĝas por iom da tempo",
		"Fantomo: ĝi kaptas vin, vi perdas vivon",
		"Timigita fantomo: manĝu ĝin, +100 (ekbrilante ĝi baldaŭ revenas)",
		"Frukto: +200, ĝi aperas de tempo al tempo · Vivo: +1 vivo",
		"Iuj muroj haltigas nur vin; nevideblaj muroj nur la fantomojn.",
		"Elirante ĉe unu rando, vi revenas ĉe la alia.",
	},
	"about": {
		"# Pri",
		"* Pacmanio — de Joop Kiefte, 2006, farita per Game Maker",
		"surbaze de la Pac-Man-ekzemplo de Mark Overmars; kelkaj bildoj",
		"de Billy McQuown (laŭ la originala helpteksto).",
		"Refarita en 2026 en Go (Ebitengine), paŝo post paŝo laŭ la originalo.",
		"",
		"Muziko kaj sonoj: el la ludo de 2006 (muziko per FluidR3, MIT).",
		"Tiparo: Noto Sans (SIL OFL). Sono malkodita per gopus.",
		"La bildoj estas la originalaj de 2006 (vidu README); la kodo: MIT.",
		"",
		"Fontkodo: github.com/LaPingvino/pacmanio",
	},
}

// drawInfo draws an info page on a panel over the title: "# " a heading,
// "* " a bold line; Enter, Esc or a click goes back.
func (a *app) drawInfo(o *game.Overlay, page string) {
	o.AddRect(game.Rect{X: 26, Y: 36, W: game.ScreenW - 52, H: game.ScreenH - 72, Color: color.RGBA{0x0b, 0x12, 0x3a, 238}})
	text, gold := color.RGBA{0xf0, 0xf4, 0xff, 255}, color.RGBA{0xff, 0xe0, 0x2a, 255}
	y := 54.0
	for _, l := range infoPages[page] {
		switch {
		case len(l) > 2 && l[:2] == "# ":
			o.Add(game.Text{S: l[2:], X: game.ScreenW / 2, Y: y, Size: 26, Weight: game.Black, Color: gold, Align: game.Centre})
			y += 40
		case len(l) > 2 && l[:2] == "* ":
			o.Add(game.Text{S: l[2:], X: 42, Y: y, Size: 14, Weight: game.Bold, Color: text})
			y += 22
		case l == "":
			y += 10
		default:
			o.Add(game.Text{S: l, X: 42, Y: y, Size: 14, Color: text})
			y += 22
		}
	}
	o.Add(game.Text{S: "Enter / Esc: reen", X: game.ScreenW / 2, Y: game.ScreenH - 62, Size: 13, Weight: game.Bold, Color: gold, Align: game.Centre})
}
