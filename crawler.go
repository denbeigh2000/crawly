package crawly

import "fmt"

type Crawler interface {
	Crawl(url string) ([]string, error)
}

type CrawlError struct {
	URL string
	Err error
}

func (e CrawlError) Error() string {
	return fmt.Sprintf("Error crawling %v (%v)", e.URL, e.Err.Error())
}
