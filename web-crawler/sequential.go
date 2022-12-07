package webcrawler

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

var fetched map[string]bool = make(map[string]bool)

// Crawl uses findLinks to recursively crawl
// pages starting with url, to a maximum of depth.
func CrawlSequential(url string, depth int) {
	if depth < 0 {
		return
	}

	// Find all links in given url
	urls, err := findLinks(url)
	if err != nil {
		return
	}

	fmt.Printf("Fetched %s\n", url)
	fetched[url] = true
	
	for _, url := range urls {
		if !fetched[url] {
			CrawlSequential(url, depth-1)
		}
	}
}

func findLinks(url string) ([]string, error) {
	// Get the page
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot get URL %s: %s", url, resp.Status)		
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse %s as HTML: %v", url, err)
	}

	return visit(nil, doc), nil
}

func visit(links []string, root *html.Node) []string {
	if root.Type == html.ElementNode && root.Data == "a" {
		for _, attr := range root.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}

	// Recursively traverse the document tree
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}

	return links
}