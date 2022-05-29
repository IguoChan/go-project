package config

import "embed"

//go:embed *.yml
var ConfigFS embed.FS
