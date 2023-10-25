package main

import (
	"fmt"
	"log"

	"hipeople.io/interview/pkg/extwo"
)

func main() {
	// Sample list of URLs to scrape
	urls := []string{
		"https://www.hipeople.io/",
		"https://github.com/",
	}

	scraper := extwo.NewCollyScrapper()

	// Scrape a single URL
	fmt.Println("Scraping a single URL...")
	data, err := scraper.Scrape(urls[0])
	if err != nil {
		log.Fatalf("Failed to scrape %s: %v", urls[0], err)
	}
	printScrapedData(data)

	// Scrape multiple URLs
	fmt.Println("\nScraping multiple URLs...")
	results, err := scraper.ScrapeMultiple(urls)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}
	for _, d := range results {
		printScrapedData(d)
	}
}

func printScrapedData(data extwo.ScrapedData) {
	fmt.Printf("URL: %s\n", data.URL)
	fmt.Printf("Title: %s\n", data.Title)
	fmt.Printf("Body: %s\n", data.Body)
	fmt.Println("Links:")
	for _, link := range data.Links {
		fmt.Printf("\t%s\n", link)
	}
	fmt.Println("----------")
}
