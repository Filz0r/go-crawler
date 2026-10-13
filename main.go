package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/filz0r/go-crawler/internal/crawler"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseUrl := args[0]

	maxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("invalid max pages provided: %s\n", err)
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("invalid max concurrency provided: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", baseUrl)

	cfg, err := crawler.NewCrawler(baseUrl, maxPages, maxConcurrency)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.StartCrawling()

	fmt.Printf("Finished crawling %s, found %d pages\n", baseUrl, len(cfg.Pages))
	for normalizedURL := range cfg.Pages {
		fmt.Printf(" - %+v\n", normalizedURL)
	}

}
