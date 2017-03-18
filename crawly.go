package crawly

import "log"

type Crawly struct {
	Concurrency int

	in chan string

	Crawler

	processor *Processor
}

func NewCrawly(start string, crawler Crawler, concurrency int) *Crawly {
	in := make(chan string)
	results := make(chan FetchResult)

	p := NewProcessor(results)

	for i := 0; i < concurrency; i++ {
		f := NewFetcher(in, crawler)

		go func(f Fetcher) {
			fetcherResults := f.Results()
			for result := range fetcherResults {
				results <- result
			}
		}(f)

		go func(f Fetcher, i int) {
			errs := f.Errors()
			for err := range errs {
				log.Printf("Error from fetcher %v: %v", err, i)
			}
		}(f, i)
	}

	c := &Crawly{
		Concurrency: concurrency,
		Crawler:     crawler,

		in:        in,
		processor: p,
	}

	return c
}

func (c *Crawly) loop() {
	defer close(c.in)

	newURLs := c.processor.New()
	for url := range newURLs {
		c.in <- url
	}
}
