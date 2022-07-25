package demo_app

import "embed"

//go:embed config.yml
var ConfigFS embed.FS
