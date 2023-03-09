package ui

import (
	"embed"
)

// Store our UI files into an embedded filesystem
//
//go:embed "html" "static"
var Files embed.FS
