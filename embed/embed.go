package embed

import "embed"

//go:embed all:schematics/*
var SchematicsFS embed.FS
