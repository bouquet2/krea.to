package converter

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// StartServer starts a development HTTP server serving the specified output directory
func StartServer(outputDir string, port string) error {
	logger := log.With().Str("dir", outputDir).Str("port", port).Logger()

	// Use the standard file server
	http.Handle("/", http.FileServer(http.Dir(outputDir)))

	addr := ":" + port
	logger.Info().Str("addr", "http://localhost:"+port).Msg("Server started")
	logger.Info().Msg("Press Ctrl+C to stop the server")

	return http.ListenAndServe(addr, nil)
}
