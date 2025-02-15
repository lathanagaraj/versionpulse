package vp

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scrapper struct {
	Url string
}

func NewScrapper(url string) *Scrapper {
	return &Scrapper{Url: url}
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

	return text, nil
}
