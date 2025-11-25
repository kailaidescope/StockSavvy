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

type ArticlePublisher struct {
	Name        string `bson:"name,omitempty"`
	HomepageURL string `bson:"homepage_url,omitempty"`
	LogoURL     string `bson:"logo_url,omitempty"`
	FaviconURL  string `bson:"favicon_url,omitempty"`
}

type ArticleInsight struct {
	Ticker             string `bson:"ticker,omitempty"`
	Sentiment          string `bson:"sentiment,omitempty"`
	SentimentReasoning string `bson:"sentiment_reasoning,omitempty"`
}

type Article struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	PolygonID   string             `bson:"polygon_id,omitempty"`
	Publisher   ArticlePublisher   `bson:"publisher,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Author      string             `bson:"author,omitempty"`
	PublishedAt primitive.DateTime `bson:"published_at,omitempty"`
	ArticleURL  string             `bson:"article_url,omitempty"`
	Tickers     []string           `bson:"tickers,omitempty"`
	ImageURL    string             `bson:"image_url,omitempty"`
	Description string             `bson:"description,omitempty"`
	Keywords    []string           `bson:"keywords,omitempty"`
	Insights    []ArticleInsight   `bson:"insights,omitempty"`
}

// InsertArticles inserts the provided articles into the "ticker_news" collection of dbName.
// It returns the number of successfully inserted documents and an error (if any).
func InsertArticles(client *mongo.Client, dbName string, articles []Article) (int, error) {
	if client == nil {
		return 0, mongo.ErrClientDisconnected
	}
	if len(articles) == 0 {
		return 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_news")

	models := make([]mongo.WriteModel, 0, len(articles))
	for _, a := range articles {
		// If polygon_id is present, use an upsert with $setOnInsert so we only insert when no document
		// with the same polygon_id exists. If polygon_id is empty, fall back to a plain insert.
		if a.PolygonID == "" {
			models = append(models, mongo.NewInsertOneModel().SetDocument(a))
		} else {
			filter := bson.M{"polygon_id": a.PolygonID}
			update := bson.M{"$setOnInsert": a}
			models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
		}
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := coll.BulkWrite(ctx, models, opts)
	if err != nil {
		// If some writes succeeded before the error, return that count along with the error.
		if res != nil {
			return int(res.InsertedCount + res.UpsertedCount), err
		}
		return 0, err
	}

	return int(res.InsertedCount + res.UpsertedCount), nil
}

// GetArticlesByTicker returns up to `limit` articles containing `ticker` in the `tickers` array,
// sorted by `published_at` descending. If limit <= 0 a default of 100 is used.
func GetArticlesByTicker(client *mongo.Client, dbName, ticker string, limit int) ([]Article, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}
	if ticker == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_news")

	filter := bson.M{"tickers": ticker}
	findOpts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []Article
	for cursor.Next(ctx) {
		var a Article
		if err := cursor.Decode(&a); err != nil {
			// skip malformed docs but continue
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return out, err
	}
	return out, nil
}

// GetArticlesByTickerOverRange returns up to `limit` articles for `ticker` where
// published_at is in [start, end). If start.IsZero() it is treated as Unix epoch.
// If end.IsZero() it is treated as time.Now(). If limit <= 0 a default of 100 is used.
func GetArticlesByTickerOverRange(client *mongo.Client, dbName, ticker string, start, end time.Time, limit int) ([]Article, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}
	if ticker == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 100
	}
	if start.IsZero() {
		start = time.Unix(0, 0).UTC()
	}
	if end.IsZero() {
		end = time.Now().UTC()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_news")

	filter := bson.M{
		"tickers": ticker,
		"published_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(start),
			"$lt":  primitive.NewDateTimeFromTime(end),
		},
	}
	findOpts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []Article
	for cursor.Next(ctx) {
		var a Article
		if err := cursor.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return out, err
	}
	return out, nil
}

// GetArticlesOverRange returns up to `limit` articles with published_at in [start, end).
// If start.IsZero() it is treated as Unix epoch. If end.IsZero() it is treated as time.Now().
// If limit <= 0 a default of 100 is used.
func GetArticlesOverRange(client *mongo.Client, dbName string, start, end time.Time, limit int) ([]Article, error) {
	if client == nil {
		return nil, mongo.ErrClientDisconnected
	}
	if limit <= 0 {
		limit = 100
	}
	if start.IsZero() {
		start = time.Unix(0, 0).UTC()
	}
	if end.IsZero() {
		end = time.Now().UTC()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	coll := client.Database(dbName).Collection("ticker_news")

	filter := bson.M{
		"published_at": bson.M{
			"$gte": primitive.NewDateTimeFromTime(start),
			"$lt":  primitive.NewDateTimeFromTime(end),
		},
	}
	findOpts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(int64(limit))

	cursor, err := coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []Article
	for cursor.Next(ctx) {
		var a Article
		if err := cursor.Decode(&a); err != nil {
			continue
		}
		out = append(out, a)
	}
	if err := cursor.Err(); err != nil {
		return out, err
	}
	return out, nil
}

// PolygonNewsToArticles converts a polygon.PolygonGetTickerNews value into a slice of mongodb Article.
func PolygonNewsToArticles(news polygon.PolygonGetTickerNews) ([]Article, error) {
	if news.Results == nil || len(*news.Results) == 0 {
		return nil, nil
	}

	results := *news.Results
	out := make([]Article, 0, len(results))

	for _, r := range results {
		var a Article
		a.ID = primitive.NewObjectID()

		if r.ID != nil {
			a.PolygonID = *r.ID
		}

		if r.Publisher != nil {
			if r.Publisher.Name != nil {
				a.Publisher.Name = *r.Publisher.Name
			}
			if r.Publisher.HomepageURL != nil {
				a.Publisher.HomepageURL = *r.Publisher.HomepageURL
			}
			if r.Publisher.LogoURL != nil {
				a.Publisher.LogoURL = *r.Publisher.LogoURL
			}
			if r.Publisher.FaviconURL != nil {
				a.Publisher.FaviconURL = *r.Publisher.FaviconURL
			}
		}

		if r.Title != nil {
			a.Title = *r.Title
		}
		if r.Author != nil {
			a.Author = *r.Author
		}
		if r.PublishedUTC != nil {
			a.PublishedAt = primitive.NewDateTimeFromTime(*r.PublishedUTC)
		}
		if r.ArticleURL != nil {
			a.ArticleURL = *r.ArticleURL
		}
		if r.Tickers != nil {
			a.Tickers = make([]string, len(*r.Tickers))
			copy(a.Tickers, *r.Tickers)
		}
		if r.ImageURL != nil {
			a.ImageURL = *r.ImageURL
		}
		if r.Description != nil {
			a.Description = *r.Description
		}
		if r.Keywords != nil {
			a.Keywords = make([]string, len(*r.Keywords))
			copy(a.Keywords, *r.Keywords)
		}
		if r.Insights != nil {
			for _, ins := range *r.Insights {
				var ai ArticleInsight
				if ins.Ticker != nil {
					ai.Ticker = *ins.Ticker
				}
				if ins.Sentiment != nil {
					ai.Sentiment = *ins.Sentiment
				}
				if ins.SentimentReasoning != nil {
					ai.SentimentReasoning = *ins.SentimentReasoning
				}
				a.Insights = append(a.Insights, ai)
			}
		}

		out = append(out, a)
	}

	return out, nil
}
