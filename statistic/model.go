package statistic

import (
	"context"

	"github.com/InsideGallery/brf.im/shorter"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/InsideGallery/core/db/mongodb"
)

type Statistic struct {
	client *mongodb.MongoClient
}

func New() (*Statistic, error) {
	db, err := mongodb.Default()
	if err != nil {
		return nil, err
	}

	return &Statistic{client: db}, nil
}

func (s *Statistic) Track(ctx context.Context, shortID string) error {
	filter := bson.D{{Key: "short_id", Value: shortID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "clicks", Value: 1}}}}

	return s.client.UpsertOne(ctx, shorter.CollectionShortURLs, update, filter)
}
