package sound

import (
	"io/fs"
	"path"
	"strings"
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

// The background preload decodes every effect, so no first play has to.
func TestPreloadDecodesEveryEffect(t *testing.T) {
	fsys := pacmanio.Files
	p := &Player{fsys: fsys}
	p.startPreload()
	<-p.done
	names, _ := fs.Glob(fsys, "pacmanio/sounds/*.wav")
	if len(names) == 0 {
		t.Fatal("no effects")
	}
	for _, n := range names {
		k := strings.TrimSuffix(path.Base(n), ".wav")
		if len(p.pcm[k]) == 0 {
			t.Errorf("%s not decoded", k)
		}
	}
	t.Logf("%d effects decoded in the background", len(names))
}
