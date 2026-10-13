package parser

import (
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func GetHeadingFromHTML(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}

	heading := doc.Find("h1")
	if heading.Length() == 0 {
		heading = doc.Find("h2")
	}
	return heading.Text()
}

func GetFirstParagraphFromHTML(html string) string {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	var paragraph *goquery.Selection
	main := doc.Find("main")
	if main.Length() == 0 {
		paragraph = doc.Find("p")
	} else {
		paragraph = main.Find("p")
	}
	return paragraph.Text()
}

func GetURLsFromHTML(html string, baseUrl *url.URL) ([]string, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	urls := make([]string, 0)
	docURLs := doc.Find("a[href]")
	if docURLs.Length() == 0 {
		return urls, nil
	}
	var parseErr error
	docURLs.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		href, found := s.Attr("href")
		if !found {
			return true
		}
		hrefURL, err := url.Parse(href)
		if err != nil {
			parseErr = err
			return false
		}
		urls = append(urls, baseUrl.ResolveReference(hrefURL).String())
		return true
	})
	if parseErr != nil {
		return nil, parseErr
	}
	return urls, nil
}

func GetImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	images := make([]string, 0)
	docImages := doc.Find("img[src]")
	if docImages.Length() == 0 {
		return images, nil
	}
	var parseErr error
	docImages.EachWithBreak(func(_ int, s *goquery.Selection) bool {
		src, found := s.Attr("src")
		if !found {
			return true
		}
		srcURL, err := url.Parse(src)
		if err != nil {
			parseErr = err
			return false
		}
		images = append(images, baseURL.ResolveReference(srcURL).String())
		return true
	})
	if parseErr != nil {
		return nil, parseErr
	}
	return images, nil
}

func ExtractPageData(html, pageURL string) (PageData, error) {
	result := PageData{}
	_url, err := url.Parse(pageURL)
	if err != nil {
		return PageData{}, err
	}
	result.URL = pageURL
	result.Heading = GetHeadingFromHTML(html)
	result.FirstParagraph = GetFirstParagraphFromHTML(html)
	outgoingLinks, err := GetURLsFromHTML(html, _url)
	if err != nil {
		return PageData{}, err
	}
	result.OutgoingLinks = outgoingLinks
	imageURLs, err := GetImagesFromHTML(html, _url)
	if err != nil {
		return PageData{}, err
	}
	result.ImageURLs = imageURLs
	return result, nil
}
