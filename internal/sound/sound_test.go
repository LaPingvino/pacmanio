package sound

import (
	"testing"

	"pacmanio"
)

// Every effect the game plays decodes.
func TestEffectsDecode(t *testing.T) {
	for _, n := range []string{"beeb", "bonus", "catch", "click", "dead", "won"} {
		pcm, err := Decode(pacmanio.Files, n)
		if err != nil || len(pcm) == 0 {
			t.Errorf("%s: %d bytes, %v", n, len(pcm), err)
		}
	}
}
