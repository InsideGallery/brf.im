package middlewares

import (
	"log/slog"

	"github.com/InsideGallery/brf.im/statistic"
	"github.com/gofiber/fiber/v2"
)

// New creates a new middleware handler
func New(s *statistic.Statistic) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		defer func() {
			shortID := c.Params("shortID")
			if shortID != "" {
				err = s.Track(c.Context(), shortID)
				if err != nil {
					slog.Error("error track redirect", "err", err)
				}
			}
		}()

		return c.Next()
	}
}
