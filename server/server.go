package server

import (
	"embed"
	_ "embed"

	"github.com/NicoClack/cryptic-stash/common"
)

var (
	//go:embed templates/*.html
	templateFS embed.FS
	//go:embed public
	publicFS embed.FS
)

var (
	TemplateFiles = common.EmbeddedDirectory{
		Path: "templates",
		FS:   templateFS,
	}
	PublicFiles = common.EmbeddedDirectory{
		Path: "public",
		FS:   publicFS,
	}
)
