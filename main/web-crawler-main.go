package main

import (
	"fmt"
	"time"

	webcrawler "edocode.me/go-concurrency-practice/web-crawler"
)

func main() {
	now := time.Now()
	// webcrawler.CrawlSequential("http://andcloud.io", 2)
	webcrawler.CrawlParallel("http://andcloud.io", 2)
	fmt.Println("Time taken:", time.Since(now))
}