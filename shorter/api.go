package shorter

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	qrcode "github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/InsideGallery/core/errors"
	"github.com/InsideGallery/core/server/webserver"
)

var (
	ErrGettingShortID error = errors.New("error getting short id from request")
	ErrInvalidPrefix  error = errors.New("error invalid prefix")
)

const (
	maxPrefixLength = 11
)

var urlLink = GetEnv("URL_LINK")

func GetEnv(name string) string {
	e := os.Getenv(name)
	if e == "" {
		return `https://brf.im`
	}

	u, err := url.Parse(e)
	if err != nil {
		return `https://brf.im`
	}

	return u.String()
}

type CreateShortURLRequest struct {
	URL    string `json:"url"`
	Prefix string `json:"prefix"`
}

func (req CreateShortURLRequest) GetPrefix() (string, error) {
	if strings.Contains(req.Prefix, "/") || len(req.Prefix) > maxPrefixLength {
		return "", ErrInvalidPrefix
	}

	return req.Prefix, nil
}

func CreateOwnerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ownerID, err := CreateOwner(c.Context())
		if err != nil {
			slog.Error("Error creating owner", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error creating owner")

			return err
		}

		requestID := c.Get("requestID")

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusCreated)

		resp := webserver.GetSuccessResponse(map[string]string{
			"owner": ownerID.Hex(),
		})

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func RemoveOwnerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner := c.Params("owner")

		id, err := primitive.ObjectIDFromHex(owner)
		if err != nil {
			slog.Error("Error decoding owner", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding owner")

			return err
		}

		err = RemoveOwner(c.Context(), id)
		if err != nil {
			slog.Error("Error removing owner", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error removing owner")

			return err
		}

		requestID := c.Get("requestID")

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusAccepted)

		resp := webserver.GetSuccessResponse(nil)

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func CreateShortURLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req CreateShortURLRequest

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			slog.Error("Error decoding request", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding request")

			return err
		}

		owner := c.Params("owner")

		id, err := primitive.ObjectIDFromHex(owner)
		if err != nil {
			slog.Error("Error decoding owner", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding owner")

			return err
		}

		prefix, err := req.GetPrefix()
		if err != nil {
			slog.Error("Error prefix is invalid", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error prefix is invalid")

			return err
		}

		shortID, err := CreateShortURL(c.Context(), prefix, req.URL, id)
		if err != nil {
			slog.Error("Error creating short url", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error creating short url")

			return err
		}

		requestID := c.Get("requestID")

		shortURL := strings.Join([]string{urlLink, "/", url.PathEscape(shortID)}, "")

		png, err := qrcode.Encode(shortURL, qrcode.Medium, 256) // nolint:mnd
		if err != nil {
			slog.Error("Error creating qr code", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error creating qr code")

			return err
		}

		sEnc := base64.StdEncoding.EncodeToString(png)

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusCreated)

		resp := webserver.GetSuccessResponse(map[string]any{
			"qrCode":    sEnc,
			"shortID":   url.PathEscape(shortID),
			"shortURL":  shortURL,
			"qrCodeURL": strings.Join([]string{urlLink, "/qr/", url.PathEscape(shortID)}, ""),
		})

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func RemoveShortURLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortID := c.Params("shortID")
		owner := c.Params("owner")

		id, err := primitive.ObjectIDFromHex(owner)
		if err != nil {
			slog.Error("Error decoding owner", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding owner")

			return err
		}

		err = RemoveShortURL(c.Context(), shortID, id)
		if err != nil {
			slog.Error("Error removing short url", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error removing short url")

			return err
		}

		requestID := c.Get("requestID")

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusAccepted)

		resp := webserver.GetSuccessResponse(nil)

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func GetShortURLsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		owner := c.Params("owner")

		id, err := primitive.ObjectIDFromHex(owner)
		if err != nil {
			slog.Error("Error decoding owner", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding owner")

			return err
		}

		urls, err := GetShortURLs(c.Context(), id)
		if err != nil {
			slog.Error("Error getting short urls", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error getting short urls")

			return err
		}

		requestID := c.Get("requestID")

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusOK)

		resp := webserver.GetSuccessResponse(map[string]any{
			"urls": urls,
		})

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func GetShortURLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortID := c.Params("shortID")
		owner := c.Params("owner")

		id, err := primitive.ObjectIDFromHex(owner)
		if err != nil {
			slog.Error("Error decoding owner", "err", err)

			c.Status(http.StatusBadRequest)
			_, err := c.WriteString("Error decoding owner")

			return err
		}

		shortURL, err := GetShortURL(c.Context(), shortID, id)
		if err != nil {
			slog.Error("Error getting short url", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error getting short url")

			return err
		}

		requestID := c.Get("requestID")

		c.Response().Header.Set("requestID", requestID)
		c.Status(http.StatusOK)

		resp := webserver.GetSuccessResponse(map[string]any{
			"url":     shortURL.URL,
			"shortID": shortURL.ShortID,
			"owner":   shortURL.Owner.Hex(),
		})

		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		_, err = c.Write(data)

		return err
	}
}

func OpenShortURLHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortID := c.Params("shortID")

		u, err := GetFullURL(c.Context(), shortID)
		if err != nil {
			slog.Error("Error getting short url", "err", err, "shortID", shortID)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error getting short url")

			return err
		}

		rawURL, err := url.Parse(u)
		if err != nil {
			slog.Error("Error parse url", "err", err, "shortID", shortID)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error parse url")

			return err
		}

		if rawURL.RawQuery != "" {
			rawURL.RawQuery = rawURL.RawQuery + "&" + c.Context().QueryArgs().String()
		} else {
			rawURL.RawQuery = c.Context().QueryArgs().String()
		}

		return c.Redirect(rawURL.String(), http.StatusPermanentRedirect)
	}
}

func GetShortURLQRCodeHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortID := c.Params("shortID")

		shortURL := strings.Join([]string{urlLink, "/", url.PathEscape(shortID)}, "")

		png, err := qrcode.Encode(shortURL, qrcode.Medium, 256) // nolint:mnd
		if err != nil {
			slog.Error("Error creating qr code", "err", err)

			c.Status(http.StatusInternalServerError)
			_, err := c.WriteString("Error creating qr code")

			return err
		}

		c.Status(http.StatusOK)
		c.Response().Header.Set("Content-Type", "image/png")

		_, err = c.Write(png)

		return err
	}
}
