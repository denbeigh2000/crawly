package crawly

import (
	"log"
	"sync"
)

type Crawly struct {
	Concurrency int

	in      chan string
	results chan FetchResult

	Crawler

	processor *Processor
}

func NewCrawly(start string, crawler Crawler, concurrency int) *Crawly {
	in := make(chan string, concurrency)
	results := make(chan FetchResult)

	c := &Crawly{
		Concurrency: concurrency,
		Crawler:     crawler,

		results: results,

		in:        in,
		processor: NewProcessor(results),
	}

	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		f := NewFetcher(in, crawler)
		go c.processFetcherResults(f, wg)
		go c.processFetcherErrors(f, i)
	}

	go c.loop()
	go func() {
		c.in <- start

		wg.Wait()
		close(in)
	}()

	return c
}

func (c *Crawly) processFetcherErrors(f Fetcher, i int) {
	errs := f.Errors()
	for err := range errs {
		log.Printf("Error from fetcher %v: %v", err, i)
	}
}

func (c *Crawly) processFetcherResults(f Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	fetcherResults := f.Results()
	for result := range fetcherResults {
		c.results <- result
	}
}

func (c *Crawly) URLs() <-chan string {
	return c.processor.Crawled()
}

func (c *Crawly) loop() {
	newURLs := c.processor.New()
	for url := range newURLs {
		c.in <- url
	}
}
