package embedded

import "embed"

//go:embed s/*
var source embed.FS

// GetSource return embed FS
func GetSource() embed.FS {
	return source
}
