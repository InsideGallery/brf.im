package shorter

import (
	"context"
	"math/rand/v2"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/InsideGallery/core/db/mongodb"
	"github.com/InsideGallery/core/errors"
	"github.com/InsideGallery/core/utils"
)

var ErrPrefixToLong error = errors.New("prefix too long")

const (
	CollectionOwner     = "owner"
	CollectionShortURLs = "short_urls"

	defaultRetries = 3
)

type OwnerModel struct {
	ID primitive.ObjectID `bson:"_id" json:"owner"`
}

type ShortURLModel struct {
	ShortID string             `bson:"short_id" json:"shortID"`
	Owner   primitive.ObjectID `bson:"owner" json:"owner"`
	URL     string             `bson:"url" json:"url"`
}

func CreateOwner(ctx context.Context) (primitive.ObjectID, error) {
	db, err := mongodb.Default()
	if err != nil {
		return primitive.ObjectID{}, err
	}
	id := primitive.NewObjectID()
	err = db.InsertOne(ctx, CollectionOwner, &OwnerModel{
		ID: id,
	})

	return id, err
}

func RemoveOwner(ctx context.Context, id primitive.ObjectID) error {
	db, err := mongodb.Default()
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "owner", Value: id}}

	_, err = db.Collection(CollectionShortURLs).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	filter = bson.D{{Key: "_id", Value: id}}

	err = db.DeleteOne(ctx, CollectionOwner, filter)
	if err != nil {
		return err
	}

	return err
}

func CreateShortURL(ctx context.Context, prefix, url string, owner primitive.ObjectID) (string, error) {
	db, err := mongodb.Default()
	if err != nil {
		return "", err
	}

	var shortID string
	var retries int

	for {
		rawShortID, err := utils.GetTinyID()
		if err != nil {
			return "", err
		}

		if retries >= defaultRetries {
			rawShortID = append(rawShortID, GetRandomChars(retries/defaultRetries)...) // every 3 retries add character
		}

		shortID = string(rawShortID)
		if prefix != "" {
			shortID = strings.Join([]string{prefix, shortID[2:]}, "-")
		}

		filter := bson.D{{Key: "short_id", Value: shortID}}

		count, err := db.CountDocuments(ctx, CollectionShortURLs, filter)
		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}

		retries++
	}
	err = db.InsertOne(ctx, CollectionShortURLs, &ShortURLModel{
		ShortID: shortID,
		Owner:   owner,
		URL:     url,
	})

	return shortID, err
}

func RemoveShortURL(ctx context.Context, shortID string, owner primitive.ObjectID) error {
	db, err := mongodb.Default()
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "short_id", Value: shortID}, {Key: "owner", Value: owner}}

	err = db.DeleteOne(ctx, CollectionShortURLs, filter)
	if err != nil {
		return err
	}

	return err
}

func GetShortURLs(ctx context.Context, owner primitive.ObjectID) ([]ShortURLModel, error) {
	shortURLModel := new(ShortURLModel)

	db, err := mongodb.Default()
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "owner", Value: owner}}
	data, err := db.Find(ctx, CollectionShortURLs, shortURLModel, filter)
	result := make([]ShortURLModel, len(data))

	for i, a := range data {
		result[i] = a.(ShortURLModel)
		result[i].ShortID = url.PathEscape(result[i].ShortID)
	}

	return result, err
}

func GetShortURL(ctx context.Context, shortID string, owner primitive.ObjectID) (*ShortURLModel, error) {
	shortURLModel := new(ShortURLModel)

	db, err := mongodb.Default()
	if err != nil {
		return shortURLModel, err
	}

	filter := bson.D{{Key: "short_id", Value: shortID}, {Key: "owner", Value: owner}}
	err = db.FindOne(ctx, CollectionShortURLs, shortURLModel, filter)
	shortURLModel.ShortID = url.PathEscape(shortURLModel.ShortID)

	return shortURLModel, err
}

func GetFullURL(ctx context.Context, shortID string) (string, error) {
	shortURLModel := new(ShortURLModel)

	db, err := mongodb.Default()
	if err != nil {
		return "", err
	}
	filter := bson.D{{Key: "short_id", Value: shortID}}
	err = db.FindOne(ctx, CollectionShortURLs, shortURLModel, filter)

	return shortURLModel.URL, err
}

var chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func GetRandomChars(n int) []byte {
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = chars[rand.IntN(len(chars))] //nolint:gosec
	}

	return result
}
