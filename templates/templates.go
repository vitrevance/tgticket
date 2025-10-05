package templates

import "embed"

//go:embed *.html
var TemplateFiles embed.FS
