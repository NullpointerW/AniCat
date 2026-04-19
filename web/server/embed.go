package webserver

import "embed"

//go:embed static
var staticFiles embed.FS
