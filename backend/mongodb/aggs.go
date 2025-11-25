package mongodb

import (
	"context"
	"financial-helper/polygon"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// InsertAggregates inserts the provided daily aggregates into the "ticker_aggregates" collection of dbName.
// It returns the number of successfully inserted documents and an error (if any).
func InsertAggregates(client *mongo.Client, dbName string, aggregates []TickerDailyAggregate) (int, error) {
	if client == nil {
		return 0, mongo.ErrClientDisconnected
	}
	if len(aggregates) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_aggregates")

	// Build a list of unique (ticker,timestamp) pairs from input and $or clauses for a single DB query.
	ors := make([]bson.M, 0, len(aggregates))
	inputMap := make(map[string]TickerDailyAggregate, len(aggregates))
	for _, a := range aggregates {
		if a.Ticker == "" {
			continue
		}
		if a.Timestamp == primitive.DateTime(0) {
			// skip invalid/missing timestamps for safety
			continue
		}
		key := fmt.Sprintf("%s_%d", a.Ticker, int64(a.Timestamp))
		if _, ok := inputMap[key]; ok {
			continue // skip duplicates in the input slice
		}
		inputMap[key] = a
		ors = append(ors, bson.M{"ticker": a.Ticker, "timestamp": a.Timestamp})
	}

	if len(ors) == 0 {
		return 0, nil
	}

	// Query the DB once to determine which pairs already exist.
	cursor, err := coll.Find(ctx, bson.M{"$or": ors}, options.Find().SetProjection(bson.M{"ticker": 1, "timestamp": 1}))
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	existing := make(map[string]struct{}, len(ors))
	var tmp struct {
		Ticker    string             `bson:"ticker"`
		Timestamp primitive.DateTime `bson:"timestamp"`
	}
	for cursor.Next(ctx) {
		if err := cursor.Decode(&tmp); err != nil {
			continue
		}
		k := fmt.Sprintf("%s_%d", tmp.Ticker, int64(tmp.Timestamp))
		existing[k] = struct{}{}
	}
	if err := cursor.Err(); err != nil {
		return 0, err
	}

	// Build insert list for aggregates that are not already present.
	toInsert := make([]interface{}, 0, len(inputMap))
	for k, a := range inputMap {
		if _, ok := existing[k]; ok {
			continue
		}
		toInsert = append(toInsert, a)
	}

	if len(toInsert) == 0 {
		return 0, nil
	}

	// Insert missing docs. Use unordered insert to maximize throughput.
	res, err := coll.InsertMany(ctx, toInsert, options.InsertMany().SetOrdered(false))
	inserted := 0
	if res != nil {
		inserted = len(res.InsertedIDs)
	}
	if err != nil {
		// If some documents were inserted before an error, return the count plus the error.
		return inserted, err
	}
	return inserted, nil
}

// GetAggregatesByTicker returns aggregates for `ticker` sorted by timestamp ascending.
// Pagination/windowing:
// - If pageSize > 0: use pagination with 1-based page; skip = (page-1)*pageSize, limit = pageSize.
// - Else if limit > 0: return up to `limit` documents (legacy behavior).
// - Else: return all matching documents.
func GetAggregatesByTicker(client *mongo.Client, dbName, ticker string, limit, page, pageSize int) ([]TickerDailyAggregate, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}
	if ticker == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_aggregates")

	filter := bson.M{"ticker": ticker}
	findOpts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	// Apply pagination or limit
	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		skip := int64((page - 1) * pageSize)
		findOpts.SetSkip(skip).SetLimit(int64(pageSize))
	} else if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]TickerDailyAggregate, 0)
	for cursor.Next(ctx) {
		var a TickerDailyAggregate
		if err := cursor.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAggregatesByTickerOverRange returns aggregates for `ticker` where timestamp is between
// `start` and `end` (inclusive). If both start and end are zero, no timestamp filter is applied.
// Pagination/windowing behavior mirrors GetAggregatesByTicker.
func GetAggregatesByTickerOverRange(client *mongo.Client, dbName, ticker string, start, end time.Time, limit, page, pageSize int) ([]TickerDailyAggregate, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}
	if ticker == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_aggregates")

	filter := bson.M{"ticker": ticker}

	// Add timestamp range filter only if start or end provided
	if !start.IsZero() || !end.IsZero() {
		rangeFilter := bson.M{}
		if !start.IsZero() {
			rangeFilter["$gte"] = primitive.NewDateTimeFromTime(start)
		}
		if !end.IsZero() {
			rangeFilter["$lte"] = primitive.NewDateTimeFromTime(end)
		}
		filter["timestamp"] = rangeFilter
	}

	findOpts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	// Apply pagination or limit
	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		skip := int64((page - 1) * pageSize)
		findOpts.SetSkip(skip).SetLimit(int64(pageSize))
	} else if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]TickerDailyAggregate, 0)
	for cursor.Next(ctx) {
		var a TickerDailyAggregate
		if err := cursor.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// GetAggregatesOverRange returns aggregates for all tickers where timestamp is between
// `start` and `end` (inclusive). If both start and end are zero, no timestamp filter is applied.
// Pagination/windowing behavior mirrors GetAggregatesByTicker.
func GetAggregatesOverRange(client *mongo.Client, dbName string, start, end time.Time, limit, page, pageSize int) ([]TickerDailyAggregate, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_aggregates")

	filter := bson.M{}

	// Add timestamp range filter only if start or end provided
	if !start.IsZero() || !end.IsZero() {
		rangeFilter := bson.M{}
		if !start.IsZero() {
			rangeFilter["$gte"] = primitive.NewDateTimeFromTime(start)
		}
		if !end.IsZero() {
			rangeFilter["$lte"] = primitive.NewDateTimeFromTime(end)
		}
		filter["timestamp"] = rangeFilter
	}

	findOpts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	// Apply pagination or limit
	if pageSize > 0 {
		if page < 1 {
			page = 1
		}
		skip := int64((page - 1) * pageSize)
		findOpts.SetSkip(skip).SetLimit(int64(pageSize))
	} else if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]TickerDailyAggregate, 0)
	for cursor.Next(ctx) {
		var a TickerDailyAggregate
		if err := cursor.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return out, nil
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
