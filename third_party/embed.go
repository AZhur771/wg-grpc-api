package third_party

import (
	"embed"
)

//go:embed openapiv2/*
var OpenAPIV2 embed.FS
