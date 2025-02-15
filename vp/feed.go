package vp

import (
	"fmt"
	"os"
	"time"

	"github.com/gorilla/feeds"
)

type ToolsFeed struct {
	ToolVersions []ToolVersion
}

func NewToolsFeed(toolVersions []ToolVersion) *ToolsFeed {
	return &ToolsFeed{
		ToolVersions: toolVersions,
	}
}

func (f *ToolsFeed) ToRss() error {
	// Create a new RSS feed
	feed := &feeds.Feed{
		Title:       "Latest Tool Versions",
		Link:        &feeds.Link{Href: "https://example.com"},
		Description: "Latest Tool Versions RSS feed",
		Created:     time.Now(),
	}

	// Convert JSON items to RSS items
	for _, toolVersion := range f.ToolVersions {

		rssItem := &feeds.Item{
			Title:       toolVersion.Tool,
			Link:        &feeds.Link{Href: toolVersion.Link},
			Description: toolVersion.Tool + " " + toolVersion.Version + " " + toolVersion.Date,
			Created:     time.Now(),
		}
		feed.Items = append(feed.Items, rssItem)
	}

	// Convert feed to RSS format
	rss, err := feed.ToJSON()
	if err != nil {
		return fmt.Errorf("error generating RSS: %v", err)
	}

	os.MkdirAll("docs", os.ModePerm)

	// Write RSS to file
	err = os.WriteFile("docs/feed.rss", []byte(rss), 0644)
	if err != nil {
		return fmt.Errorf("error writing RSS file: %v", err)
	}

	fmt.Println("RSS feed successfully generated as feed.rss")
	return nil
}
