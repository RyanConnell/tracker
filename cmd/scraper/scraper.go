package main

import (
	"fmt"
	"log"
	"os"

	"tracker/trackable/show"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log.Printf("starting scraper")
	if err := show.ScrapeAll(); err != nil {
		return fmt.Errorf("scrape error: %w", err)
	}

	return nil
}
