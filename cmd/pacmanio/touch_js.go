//go:build js

package main

import (
	"strings"
	"syscall/js"
)

// showPadsFromURL: "?pads" in the page's address shows the touch pads from the
// start (for checking the layout without a touch screen).
func showPadsFromURL() bool {
	return strings.Contains(js.Global().Get("location").Get("search").String(), "pads")
}

// pollPhoneKeys takes what the page queued in pacmanioTyped (strings: one
// character each, "\b" Backspace, "\r" Enter) and empties it.
func pollPhoneKeys() phoneKeys {
	var p phoneKeys
	q := js.Global().Get("pacmanioTyped")
	if q.IsUndefined() || q.Length() == 0 {
		return p
	}
	for i := 0; i < q.Length(); i++ {
		s := q.Index(i).String()
		if s == "\r" {
			p.enter = true
		} else {
			p.chars = append(p.chars, []rune(s)...)
		}
	}
	q.Call("splice", 0, q.Length())
	return p
}

// roomFromURL: "?room=name" starts in that room (for checking a room without
// playing up to it).
func roomFromURL() string {
	q := js.Global().Get("location").Get("search").String()
	if i := strings.Index(q, "room="); i >= 0 {
		r := q[i+5:]
		if j := strings.IndexAny(r, "&#"); j >= 0 {
			r = r[:j]
		}
		return r
	}
	return ""
}
