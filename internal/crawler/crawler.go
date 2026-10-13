package crawler

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/filz0r/go-crawler/internal/getter"
	"github.com/filz0r/go-crawler/internal/normalizing"
	"github.com/filz0r/go-crawler/internal/parser"
)

type Config struct {
	Pages              map[string]parser.PageData
	BaseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func NewCrawler(baseUrl string, maxPages, maxConcurrency int) (*Config, error) {
	if maxConcurrency <= 0 {
		return nil, fmt.Errorf("maxConcurrency must be positive: %d", maxConcurrency)
	} else if maxConcurrency > 32 {
		return nil, fmt.Errorf("maxConcurrency must be <= 32: %d", maxConcurrency)
	}
	_baseURL, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	result := &Config{
		Pages:              make(map[string]parser.PageData),
		BaseURL:            _baseURL,
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		mu:                 &sync.Mutex{},
		maxPages:           maxPages,
	}
	return result, nil
}

func (c *Config) shouldCrawl(normalizedURL string) (shouldCrawl bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Pages) >= c.maxPages {
		return false
	}

	_, seen := c.Pages[normalizedURL]
	if seen {
		return false
	}
	c.Pages[normalizedURL] = parser.PageData{}
	return true
}

func (c *Config) crawlPage(rawCurrentURL string) {
	defer c.wg.Done()
	defer func() { <-c.concurrencyControl }()
	c.concurrencyControl <- struct{}{}

	c.mu.Lock()
	atLimit := len(c.Pages) >= c.maxPages
	c.mu.Unlock()
	if atLimit {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error parsing current URL %q: %v\n", rawCurrentURL, err)
		return
	}
	if currentURL.Hostname() != c.BaseURL.Hostname() {
		fmt.Printf("Skipping external URL %s\n", rawCurrentURL)
		return
	}
	key, err := normalizing.NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error normalizing URL %q: %v\n", rawCurrentURL, err)
		return
	}

	if !c.shouldCrawl(key) {
		return
	}

	body, err := getter.GetHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error Crawling %q: %v\n", rawCurrentURL, err)
		return
	}

	page, err := parser.ExtractPageData(body, rawCurrentURL)
	if err != nil {
		fmt.Printf("Error extracting page data from %q: %v\n", rawCurrentURL, err)
		return
	}
	c.mu.Lock()
	c.Pages[key] = page
	c.mu.Unlock()

	fmt.Printf("Crawled %q\n", rawCurrentURL)

	for _, link := range page.OutgoingLinks {
		c.wg.Add(1)
		go c.crawlPage(link)
	}
}

func (c *Config) StartCrawling() {
	c.wg.Add(1)
	go c.crawlPage(c.BaseURL.String())
	c.wg.Wait()
}
