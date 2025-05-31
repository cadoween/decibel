package iox

import (
	"io"

	"github.com/rs/zerolog"
)

// Close cleanly closes an io.Closer, logging any errors encountered using
// zerolog.
func Close(closer io.Closer, logger ...*zerolog.Logger) {
	if closer == nil {
		return
	}

	if err := closer.Close(); err != nil {
		if len(logger) > 0 {
			logger[0].Error().Err(err).Msg("iox: failed to Close")
		}
	}
}
