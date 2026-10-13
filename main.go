package main

import (
	"fmt"
	"os"

	"github.com/filz0r/go-crawler/internal/crawler"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	baseUrl := args[0]
	fmt.Printf("starting crawl of: %s\n", baseUrl)

	cfg, err := crawler.NewCrawler(baseUrl, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg.StartCrawling()

	fmt.Printf("Finished crawling %s, found %d pages\n", baseUrl, len(cfg.Pages))
	for _, page := range cfg.Pages {
		fmt.Printf(" - %+v\n", page)
	}

}
