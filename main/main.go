package main

import (
	"fmt"
	"log"
	"os"
	"time"

	imageprocessing "edocode.me/go-concurrency-practice/image-processing"
	webcrawler "edocode.me/go-concurrency-practice/web-crawler"
)

func main() {
	imageProcessingMain()
}

// Usage: go run main.go ../image-processing/imgs
// save output to ./thumbnail
func imageProcessingMain() {
	if len(os.Args) < 2 {
		log.Fatal("Directory path of images required.")
	}

	start := time.Now()

	// err := imageprocessing.ProcessSequential(os.Args[1])
	err := imageprocessing.ProcessParallel(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Time taken: %s\n", time.Since(start))
}

// Usage: go run main.go
func webCrawlerMain() {
	now := time.Now()
	// webcrawler.CrawlSequential("http://andcloud.io", 2)
	webcrawler.CrawlParallel("http://andcloud.io", 2)
	fmt.Println("Time taken:", time.Since(now))
}