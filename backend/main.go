package main

import (
	"financial-helper/scraper"
	"financial-helper/server"
	"flag"
	"log"
)

func main() {
	runScraperFlag := flag.String("scrape", "", "Runs the scraper.")
	flag.Parse()

	if *runScraperFlag != "" {
		if *runScraperFlag == "aggs" {
			runAggsScraper()
		} else if *runScraperFlag == "news" {
			runNewsScraper()
		}
	} else {
		runServer()
	}
}

func runNewsScraper() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start scraper:", err)
	}

	scraper.ScrapeTickersNewsFromJSON("./scraper/article_instructions.json")
}

func runAggsScraper() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start scraper:", err)
	}

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
