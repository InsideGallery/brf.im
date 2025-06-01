package embedded

import "embed"

//go:embed default/*
var template embed.FS

// GetTemplate return embed FS
func GetTemplate() embed.FS {
	return template
}
