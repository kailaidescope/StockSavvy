package mongodb

import (
	"context"
	"financial-helper/polygon"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertHistory inserts the provided daily aggregates into the "ticker_aggregates" collection of dbName.
// It returns the number of successfully inserted documents and an error (if any).
func InsertHistory(client *mongo.Client, dbName string, aggregates []TickerDailyAggregate) (int, error) {
	if client == nil {
		return 0, mongo.ErrClientDisconnected
	}
	if len(aggregates) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_aggregates")

	models := make([]mongo.WriteModel, 0, len(aggregates))
	for _, a := range aggregates {
		// Use the composite (ticker, timestamp) as the upsert key to match the index `(ticker_timestamp)`.
		filter := bson.M{
			"ticker":    a.Ticker,
			"timestamp": a.Timestamp,
		}
		// Only insert when no document with the same (ticker, timestamp) exists.
		update := bson.M{"$setOnInsert": a}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		if res != nil {
			return int(res.InsertedCount + res.UpsertedCount), err
		}
		return 0, err
	}

	return int(res.InsertedCount + res.UpsertedCount), nil
}

type TickerDailyAggregate struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Ticker       string             `bson:"ticker,omitempty"`
	Volume       float64            `bson:"volume,omitempty"`
	VWAP         float64            `bson:"vwap,omitempty"`
	Open         float64            `bson:"open,omitempty"`
	Close        float64            `bson:"close,omitempty"`
	High         float64            `bson:"high,omitempty"`
	Low          float64            `bson:"low,omitempty"`
	Timestamp    primitive.DateTime `bson:"timestamp,omitempty"`
	Transactions int                `bson:"transactions,omitempty"`
	OTC          bool               `bson:"otc,omitempty"`
}

// Convert a PolygonGetTickerHistoryResponse into a slice of TickerDailyAggregate.
func PolygonHistoryToAggs(news polygon.PolygonGetTickerHistoryResponse) ([]TickerDailyAggregate, error) {
	if news.Results == nil || len(*news.Results) == 0 {
		return nil, nil
	}
	results := *news.Results
	out := make([]TickerDailyAggregate, 0, len(results))

	for _, r := range results {
		var a TickerDailyAggregate
		a.ID = primitive.NewObjectID()
		if news.Ticker != nil {
			a.Ticker = *news.Ticker
		}
		if r.Volume != nil {
			a.Volume = *r.Volume
		}
		if r.VWAP != nil {
			a.VWAP = *r.VWAP
		}
		if r.Open != nil {
			a.Open = *r.Open
		}
		if r.Close != nil {
			a.Close = *r.Close
		}
		if r.High != nil {
			a.High = *r.High
		}
		if r.Low != nil {
			a.Low = *r.Low
		}
		if r.Timestamp != nil {
			a.Timestamp = primitive.NewDateTimeFromTime(time.Unix(0, *r.Timestamp*int64(time.Millisecond)))
		}
		if r.Transactions != nil {
			a.Transactions = *r.Transactions
		}
		if r.OTC != nil {
			a.OTC = *r.OTC
		}
		out = append(out, a)
	}

	return out, nil
}
