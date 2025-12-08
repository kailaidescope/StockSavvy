package main

import (
	"context"
	"financial-helper/scraper"
	"financial-helper/server"
	"flag"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	runScraperFlag := flag.String("scrape", "", "Runs the scraper.")
	flag.Parse()

	if *runScraperFlag != "" {
		switch *runScraperFlag {
		case "aggs":
			runAggsScraper()
		case "news":
			runNewsScraper()
		case "sentiment":
			runSentimentUpdater()
		default:
			log.Printf("%s is not a valid scraping option", *runScraperFlag)
		}

	} else {
		runServer()
	}
}

func runSentimentUpdater() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start sentiment updater:", err)
	}
	defer scraper.MongoDisconnect(context.Background())

	queryParam := bson.D{{Key: "insights", Value: bson.M{"$exists": false}}}

	_, err = scraper.PaginateAllArticles(1, 2, queryParam)
	if err != nil {
		log.Println("Error paginating articles:", err)
	}
}

func runNewsScraper() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start scraper:", err)
	}
	defer scraper.MongoDisconnect(context.Background())

	scraper.ScrapeTickersNewsFromJSON("./scraper/article_instructions.json")
}

func runAggsScraper() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start scraper:", err)
	}
	defer scraper.MongoDisconnect(context.Background())

	scraper.ScrapeTickersAggregatesFromJSON("./scraper/aggs_instructions.json")
}

func runServer() {
	gin_server, err := server.GetNewServer()
	if err != nil {
		log.Fatal("Could not get the server object: ", err)
	}

	err = gin_server.Router.Run(":3333")
	if err != nil {
		log.Fatal("Could not start the server: ", err)
	}

	gin_server.GeminiClient.Close()
}
