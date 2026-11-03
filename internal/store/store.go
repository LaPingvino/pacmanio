// Package store keeps the highscore table between games: in the browser's
// localStorage on the web, in a JSON file in the user's config directory on
// the desktop (the original kept it in its own data).
package store

import (
	"encoding/json"
	"strconv"
	"strings"

	"pacmanio/internal/game"
)

const key = "pacmanio-highscores"

// LoadHighscores returns the saved table, or the cleared one.
func LoadHighscores() game.Highscores {
	h := game.ClearedHighscores()
	if data, ok := read(); ok {
		var saved game.Highscores
		if json.Unmarshal(data, &saved) == nil {
			h = saved
		}
	}
	return h
}

// SaveHighscores stores the table.
func SaveHighscores(h game.Highscores) {
	if data, err := json.Marshal(h); err == nil {
		write(data)
	}
}

// Progress is the furthest room reached (its index in game.RoomOrder; 0 none):
// the title's Continue ("resume wherever you got, like with Pits").
func Progress() int {
	data, ok := readKey("pacmanio-progress")
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(data)))
	return n
}

// Reached records room index n if it is further than before.
func Reached(n int) {
	if n > Progress() {
		writeKey("pacmanio-progress", []byte(strconv.Itoa(n)))
	}
}

// Muted is the sound button's setting (off: false).
func Muted() bool {
	data, ok := readKey("pacmanio-muted")
	return ok && strings.TrimSpace(string(data)) == "1"
}

// SetMuted keeps the sound button's setting.
func SetMuted(m bool) {
	v := "0"
	if m {
		v = "1"
	}
	writeKey("pacmanio-muted", []byte(v))
}

// MusicOff is the music setting (on: false).
func MusicOff() bool {
	data, ok := readKey("pacmanio-music-off")
	return ok && strings.TrimSpace(string(data)) == "1"
}

// SetMusicOff keeps the music setting.
func SetMusicOff(off bool) {
	v := "0"
	if off {
		v = "1"
	}
	writeKey("pacmanio-music-off", []byte(v))
}
