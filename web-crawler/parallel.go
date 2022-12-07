package webcrawler

import "fmt"

type result struct {
	url string
	depth int
	urlsContained []string
	err error
}

func CrawlParallel(url string, depth int) {
	resultsCh := make(chan *result)
	defer close(resultsCh)

	fetch := func (url string, depth int) {
		urls, err := findLinks(url)
		resultsCh <- &result{url, depth, urls, err}
	}

	// Find all links in given url
	go fetch(url, depth)
	fetched[url] = true

	for fetching := 1; fetching > 0; fetching-- {
		res := <- resultsCh

		if res.err != nil {
			continue
		}

		fmt.Printf("Fetched %s\n", res.url)
		if res.depth > 0 {
			for _, urlContained := range res.urlsContained {
				if !fetched[urlContained] {
					fetching++
					go fetch(urlContained, res.depth-1)
					fetched[urlContained] = true
				}
			}
		}
	}
		
}

