package main

import (
	"financial-helper/scraper"
	"financial-helper/server"
	"flag"
	"log"
)

func main() {
	runScraperFlag := flag.Bool("scrape", false, "Runs the scraper.")
	flag.Parse()

	if *runScraperFlag {
		runScraper()
	} else {
		runServer()
	}
}

func runScraper() {
	scraper, err := scraper.New()
	if err != nil {
		log.Fatal("Failed to start scraper:", err)
	}

	scraper.ScrapeTickersFromJSON("./scraper/instructions.json")
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
