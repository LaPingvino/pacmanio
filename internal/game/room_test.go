package game

import "testing"

// muro's room-start table (pacmanio/objects/muro.gml), case by case:
// neighbours left, up, right, down → frame.
func TestWallFrames(t *testing.T) {
	cases := map[[4]bool]int{
		{}:                          0,
		{false, false, false, true}: 2,
		{false, false, true, false}: 4,
		{false, false, true, true}:  13,
		{false, true, false, false}: 1,
		{false, true, false, true}:  14,
		{false, true, true, false}:  10,
		{false, true, true, true}:   9,
		{true, false, false, false}: 3,
		{true, false, false, true}:  12,
		{true, false, true, false}:  15,
		{true, false, true, true}:   8,
		{true, true, false, false}:  11,
		{true, true, false, true}:   7,
		{true, true, true, false}:   6,
		{true, true, true, true}:    5,
	}
	for n, want := range cases {
		r := &Room{Region: 1, Things: []*Thing{{Object: "muro", X: 64, Y: 64}}}
		for i, d := range [4][2]int{{-32, 0}, {0, -32}, {32, 0}, {0, 32}} {
			if n[i] {
				r.Things = append(r.Things, &Thing{Object: "murox", X: 64 + d[0], Y: 64 + d[1]})
			}
		}
		r.autotile()
		if got := r.Things[0].Frame; got != want {
			t.Errorf("neighbours %v: frame %d, want %d", n, got, want)
		}
	}
}
