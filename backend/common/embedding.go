package common

import "embed"

type EmbeddedDirectory struct {
	Path string
	FS   embed.FS
}
