package tap

import (
	"embed"
	"io/fs"
)

const (
	EmbeddedTapWebSkillRoot   = "skills/tap-web"
	EmbeddedTapWebSkillConfig = EmbeddedTapWebSkillRoot + "/SKILL.md"
)

//go:embed skills/tap-web/*
//go:embed skills/tap-web/references/*
var embeddedTapWebSkillFS embed.FS

func EmbeddedTapWebSkillFS() fs.FS {
	return embeddedTapWebSkillFS
}
