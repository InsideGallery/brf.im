package pages

import (
	t "html/template"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/InsideGallery/core/server/template"
)

// Page describe common page
type Page struct {
	Name           string                 `bson:"name"`
	Image          string                 `bson:"image"`
	Title          string                 `bson:"title"`
	Description    string                 `bson:"description"`
	Keywords       string                 `bson:"keywords"`
	TextHead       t.HTML                 `bson:"text_head"`
	Text           t.HTML                 `bson:"text"`
	Header         t.HTML                 `bson:"header"`
	Resources      t.HTML                 `bson:"resources"`
	Additional     map[string]interface{} `bson:"additional"`
	AdditionalJSON t.HTML                 `bson:"-"`
}

// NewPage return new page
func NewPage(title, description, keywords string, header, resources t.HTML, textHead t.HTML, text t.HTML) Page {
	return Page{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		Header:      header,
		Resources:   resources,
		Text:        text,
		TextHead:    textHead,
		Additional:  make(map[string]interface{}),
	}
}

// Add add additional data to page
func (p *Page) Add(key string, value interface{}) {
	p.Additional[key] = value
}

func PageHandler(name string, tmpl *template.Engine) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Set("Content-Type", "text/html; charset=utf-8")

		pg := NewPage("Main | Brief I am", ``, ``, ``, ``, ``, ``)
		for key, value := range pg.Additional {
			switch v := value.(type) {
			case map[string]interface{}:
				for k, val := range v {
					if res, ok := val.(string); ok {
						v[k] = t.HTML(res) // nolint:gosec
					}
				}
			case string:
				pg.Additional[key] = t.HTML(v) // nolint:gosec
			}
		}

		res, err := tmpl.Execute(name, pg)
		if err != nil {
			slog.Error("Error parsing response", "err", err)
			return err
		}

		_, err = c.Write(res)
		if err != nil {
			slog.Error("Error sending data", "err", err)
			return err
		}

		return nil
	}
}

func NotFound(tmpl *template.Engine, w http.ResponseWriter) {
	res, err := tmpl.Execute("404", NewPage("404 | Brief I am", ``, ``, ``, ``, ``, ``))
	if err != nil {
		slog.Error("Error parsing response", "err", err)
		return
	}

	w.WriteHeader(http.StatusNotFound)

	_, err = w.Write(res)
	if err != nil {
		slog.Error("Error sending data", "err", err)
		return
	}
}
