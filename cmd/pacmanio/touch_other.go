//go:build !js

package main

import "os"

func showPadsFromURL() bool { return false }

func pollPhoneKeys() phoneKeys { return phoneKeys{} }

// roomFromURL: on the desktop, PACMANIO_ROOM=name starts in that room.
func roomFromURL() string { return os.Getenv("PACMANIO_ROOM") }
