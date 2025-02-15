package main

import (
	"log"
	"versionpulse/vp"
)

func main() {
	checkToolVersions()
}

func checkToolVersions() {
	tools, err := vp.Load()
	if err != nil {
		log.Fatalf("Error loading tools: %v", err)
		return
	}

	var toolVersions []vp.ToolVersion

	for _, tool := range tools.Tools {

		siteContent, err := vp.NewScrapper(tool.URL).Scrape()
		if err != nil {
			log.Printf("Error fetching %s: %v", tool.URL, err)
			continue
		}

		toolVersion, err := vp.NewChecker(tool.Name, siteContent).CheckVersion()
		if err != nil {
			log.Printf("Error extracting JSON object: %v", err)
			continue
		}
		toolVersion.Link = tool.URL
		toolVersions = append(toolVersions, *toolVersion)

		log.Print("result " + toolVersion.Tool + " " + toolVersion.Version + " " + toolVersion.Date + " " + toolVersion.Link)
	}

	if err := vp.NewToolsFeed(toolVersions).ToRss(); err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

}
