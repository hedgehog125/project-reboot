package server

import (
	"embed"

	"github.com/NicoClack/cryptic-stash/backend/common"
)

var (
	//go:embed templates/*.html
	templateFS embed.FS
	// Go excludes underscored files unless we use "all:"
	//go:embed all:public
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
