package main

import (
	"context"
	"log/slog"

	_ "github.com/InsideGallery/core/fastlog/handlers/stderr"

	"github.com/InsideGallery/brf.im/handler"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/InsideGallery/core/app"
	"github.com/InsideGallery/core/db/mongodb"
	"github.com/InsideGallery/core/fastlog/metrics"
	"github.com/InsideGallery/core/server/instance"
	"github.com/InsideGallery/core/server/profiler"
)

func main() {
	ctx := context.Background()

	app.WebMain(ctx, ":8080", "brf.im", func(
		ctx context.Context,
		app *fiber.App,
		_ *metrics.OTLPMetric,
	) error {
		mongoClient, err := mongodb.Default()
		if err != nil {
			return err
		}

		hl, err := handler.NewHandler(ctx, app, mongoClient)
		if err != nil {
			return err
		}

		profiler.AddHealthCheck(func() error {
			return mongoClient.Ping(ctx, readpref.SecondaryPreferred())
		})

		slog.Info("Instance ready", "id", instance.GetShortInstanceID())

		return hl.Run()
	})
}
