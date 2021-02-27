// Package web embeds the data. This has to be done here as the path has to
// be relative to the package, and if done in frontend, it forms a weird tree
package web

import (
	"embed"
	"io/fs"
)

//go:embed public
var Static embed.FS

//go:embed templates
var templates embed.FS

var Templates, _ = fs.Sub(templates, "templates")
