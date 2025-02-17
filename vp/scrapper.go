package vp

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scrapper struct {
	Url   string
	Limit int
}

func NewScrapper(url string, limit int) *Scrapper {
	return &Scrapper{Url: url, Limit: limit}
}

// Extracts text from the url passed.
// It only works for static pages. Does not support Dynamic pages with javascript
func (s *Scrapper) Scrape() (string, error) {
	resp, err := http.Get(s.Url)
	if err != nil {
		log.Printf("Error fetching %s: %v", s.Url, err)
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing HTML for %s: %v", s.Url, err)
		return "", err
	}

	text := doc.Text()
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.TrimSpace(text)

	return truncate(text, s.Limit), nil
}

func truncate(text string, limit int) string {
	if len(text) > limit {
		return text[:limit]
	}
	return text
}
