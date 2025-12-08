package scraper

import (
	"context"
	"errors"
	"financial-helper/mongodb"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (scraper *Scraper) PaginateAllArticles(page, pageSize int, matches ...bson.D) (*[]mongodb.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Add defaults and limits to page [1,inf) and pageSize [1,500]
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	if pageSize > 500 {
		pageSize = 500
	}

	matchStages := bson.A{}
	for _, queryParam := range matches {
		log.Printf("Adding query param: %v", queryParam)
		matchStages = append(matchStages, bson.D{{Key: "$match", Value: queryParam}})
	}

	metadataQuery := bson.A{bson.D{{Key: "$count", Value: "totalCount"}}}
	metadataQuery = append(matchStages, metadataQuery...)
	dataQuery := bson.A{bson.D{{Key: "$skip", Value: ((page - 1) * pageSize)}}, bson.D{{Key: "$limit", Value: pageSize}}}
	dataQuery = append(matchStages, dataQuery...)

	log.Printf("Full query params: %v", dataQuery)

	facetStage := bson.D{{Key: "$facet", Value: bson.D{{Key: "metadata", Value: metadataQuery}, {Key: "data", Value: dataQuery}}}}

	allArticlesReponse, err := scraper.mongoArticlesCollection.Aggregate(ctx, mongo.Pipeline{facetStage})
	if err != nil {
		return nil, errors.Join(errors.New("failed to retrieve paginated data from mongodb"), err)
	}

	// Decode response
	var articlesFacetResult []struct {
		Metadata []struct {
			TotalCount int64 `bson:"totalCount"`
		} `bson:"metadata"`
		Data []mongodb.Article `bson:"data"`
	}
	if err = allArticlesReponse.All(ctx, &articlesFacetResult); err != nil {
		return nil, errors.Join(errors.New("failed to decode paginated response into list of articles"), err)
	}

	// Check results
	if len(articlesFacetResult) < 1 {
		return nil, errors.New("no response found")
	}
	articlesResult := articlesFacetResult[0]

	log.Printf("Retrieved articles:\n%+v\n", articlesResult)

	return &articlesResult.Data, nil
}
