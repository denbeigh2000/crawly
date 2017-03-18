package main

import (
	"github.com/denbeigh2000/crawly"

	"fmt"
	"log"
)

type FakeCrawler map[string][]string

func (f FakeCrawler) Crawl(url string) ([]string, error) {
	if v, ok := f[url]; ok {
		return v, nil
	} else {
		return nil, crawly.CrawlError{URL: url, Err: fmt.Errorf("404 Not Found")}
	}
}

func main() {
	crawler := FakeCrawler{
		"http://www.somecompany.com": []string{
			"http://www.somecompany.com",
			"http://www.somecompany.com/about",
			"http://www.somecompany.com/jobs",
			"http://www.somecompany.com/login",
		},
		"http://www.somecompany.com/about": []string{
			"http://www.somecompany.com",
			"http://www.somecompany.com/jobs",
			"http://www.somecompany.com/foo",
		},
		"http://www.somecompany.com/jobs": []string{
			"http://www.somecompany.com",
		},
		"http://www.somecompany.com/foo": []string{
			"http://www.somecompany.com",
			"http://www.somecompany.com/login",
			"http://www.somecompany.com/bar",
			"http://www.somecompany.com/baz",
		},
		"http://www.somecompany.com/bar": []string{
			"http://www.somecompany.com/bar",
		},
		"http://www.somecompany.com/baz": []string{
			"http://www.somecompany.com/bar",
		},
		"http://www.somecompany.com/login": []string{
			"http://www.somecompany.com/bar",
		},
	}
	c := crawly.NewCrawly("http://www.somecompany.com", crawler, 1)

	urls := c.URLs()
	for url := range urls {
		log.Printf("Crawled URL (%v)\n", url)
	}

	log.Printf("Finished crawling")
}
