package web

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

var Files, _ = fs.Sub(dist, "dist")
