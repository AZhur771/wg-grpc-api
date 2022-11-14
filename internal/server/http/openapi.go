package httpserver

import (
	"io/fs"
	"mime"
	"net/http"

	"github.com/AZhur771/wg-grpc-api/third_party"
)

// GetOpenAPIHandler serves an OpenAPI UI.
func GetOpenAPIHandler() (http.Handler, error) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPIV2, "openapiv2")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(subFS)), nil
}
