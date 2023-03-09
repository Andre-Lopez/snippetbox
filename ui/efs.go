package ui

import (
	"embed"
)

// Store our UI files into an embedded filesystem
//
//go:embed "static"
var Static embed.FS

//go:embed "html"
var Templates embed.FS
