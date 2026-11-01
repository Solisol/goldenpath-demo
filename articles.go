// Package goldenpathdemo holds the golden-path articles, embedded so the service
// ships as one binary. In a real platform they might come from a wiki or a git repo
// instead (see GOLDENPATH_CONTENT_DIR in cmd/goldenpath-api).
package goldenpathdemo

import "embed"

//go:embed articles/*.md
var Articles embed.FS
