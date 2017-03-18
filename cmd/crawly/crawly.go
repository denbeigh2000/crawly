package main

import (
	"github.com/denbeigh2000/crawly"

	"flag"
	"fmt"
	"log"
)

type FakeCrawler map[string][]string

var (
	concurrency = flag.Int("concurrency", 1, "Number of HTTP Fetching routines to intstantiate")
	start       = flag.String("start-at", "http://www.somecompany.com", "The URL to start crawling from")
)

var defaultCrawler = FakeCrawler{
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

func (f FakeCrawler) Crawl(url string) ([]string, error) {
	if v, ok := f[url]; ok {
		return v, nil
	} else {
		return nil, crawly.CrawlError{URL: url, Err: fmt.Errorf("404 Not Found")}
	}
}

func init() {
	flag.Parse()
}

func main() {
	if *concurrency <= 0 {
		log.Panic("Concurrency must be at least 1, got %v", *concurrency)
	}
	if *start == "" {
		log.Panic("Must ptrovide a start URL", *concurrency)
	}

	c := crawly.NewCrawly("http://www.somecompany.com", defaultCrawler, *concurrency)

	urls := c.URLs()
	for url := range urls {
		log.Printf("Crawled URL (%v)\n", url)
	}

	log.Printf("Finished crawling")
}
