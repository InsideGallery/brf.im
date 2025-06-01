package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/InsideGallery/brf.im/handler/middlewares"
	"github.com/InsideGallery/brf.im/handler/pages"
	embedded "github.com/InsideGallery/brf.im/resources"
	"github.com/InsideGallery/brf.im/shorter"
	"github.com/InsideGallery/brf.im/statistic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/InsideGallery/core/db/mongodb"
	"github.com/InsideGallery/core/server/template"
)

// Handler describe handler
type Handler struct {
	*template.Engine
	ctx         context.Context
	app         *fiber.App
	mongoClient *mongodb.MongoClient
}

// NewHandler return new handler
func NewHandler(ctx context.Context, app *fiber.App, mongoClient *mongodb.MongoClient) (*Handler, error) {
	h := &Handler{
		Engine:      template.NewEngine(),
		ctx:         ctx,
		mongoClient: mongoClient,
		app:         app,
	}

	return h, nil
}

func (h *Handler) Run() error {
	// middleware := webserver.NewMiddleware(
	//	 middlewares.RecoverFiber,
	// )
	st, err := statistic.New()
	if err != nil {
		return err
	}

	h.app.Use(
		cors.New(),
		recover.New(recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(_ *fiber.Ctx, e interface{}) {
				slog.Default().Error("Recovered panic", "err", e)
			},
		}),
		middlewares.New(st),
	)

	h.app.Get("/", pages.PageHandler("main", h.Engine))
	h.app.Get("/:shortID", shorter.OpenShortURLHandler())
	h.app.Get("/qr/:shortID", shorter.GetShortURLQRCodeHandler())
	h.app.Post("/owner", shorter.CreateOwnerHandler())
	h.app.Delete("/owner/:owner", shorter.RemoveOwnerHandler())
	h.app.Post("/owner/:owner/url", shorter.CreateShortURLHandler())
	h.app.Get("/owner/:owner/url", shorter.GetShortURLHandler())
	h.app.Delete("/owner/:owner/url/:shortID", shorter.RemoveShortURLHandler())
	h.app.Get("/owner/:owner/url/:shortID", shorter.RemoveShortURLHandler())
	h.app.Use("/s", filesystem.New(filesystem.Config{
		Root:       http.FS(embedded.GetSource()),
		PathPrefix: "s",
		Browse:     true,
	}))

	tmpl, err := template.NewTemplateBySource(embedded.GetTemplate(), "main", "default/index.html")
	if err != nil {
		return err
	}

	h.Add(tmpl)

	return nil
}

// ErrorHandler default error handler
func (h *Handler) ErrorHandler(status int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Error("ErrorHandler", "method", r.Method, "url", r.URL.String())
		http.Error(w, "Error during load", status)
	}
}
