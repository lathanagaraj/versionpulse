package main

import (
	"log"
	"sync"
	"time"
	"versionpulse/vp"
)

const WebContentLenght int = 25000
const maxWorkers = 1 // Limit the number of concurrent goroutines

func main() {
	start := time.Now()
	checkToolVersions()
	elapsed := time.Since(start)
	log.Printf("Time taken: %s\n", elapsed)
}

func checkToolVersions() {
	tools, err := vp.Load()
	if err != nil {
		log.Fatalf("Error loading tools: %v", err)
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to synchronize access to toolVersions

	var toolVersions []vp.ToolVersion

	semaphore := make(chan struct{}, maxWorkers) // Buffered channel to limit concurrency

	for _, tool := range tools.Tools {

		wg.Add(1)
		semaphore <- struct{}{} // Acquire a slot

		go func(tool vp.Tool) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the slot
			siteContent, err := vp.NewScrapper(tool.URL, WebContentLenght).Scrape()
			if err != nil {
				log.Printf("Error fetching %s: %v", tool.URL, err)
				return
			}

			toolVersion, err := vp.NewChecker(tool.ID, tool.Name, siteContent).CheckVersion()
			if err != nil {
				log.Printf("Error extracting JSON object: %v", err)
				return
			}
			//time.Sleep(10 * time.Second)

			toolVersion.Link = tool.URL
			mu.Lock()
			toolVersions = append(toolVersions, *toolVersion)
			mu.Unlock()

			log.Print("result " + toolVersion.Tool + " " + toolVersion.Version + " " + toolVersion.Date + " " + toolVersion.Link)
		}(tool)
	}

	wg.Wait()
	if err := vp.NewToolsFeed(toolVersions).ToRss(); err != nil {
		log.Fatalf("Error generating RSS feed: %v", err)
	}

	log.Println("All tools processed")

}
