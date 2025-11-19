package mongodb

import (
	"financial-helper/polygon"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
