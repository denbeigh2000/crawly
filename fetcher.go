package crawly

type FetchResult struct {
	src     string
	results []string
}

type Fetcher struct {
	in  <-chan string
	out chan FetchResult

	errors chan error

	Crawler
}

func NewFetcher(urls <-chan string, c Crawler) Fetcher {
	f := Fetcher{
		in:      urls,
		out:     make(chan FetchResult),
		errors:  make(chan error),
		Crawler: c,
	}

	go f.loop()
	return f
}

func (f Fetcher) loop() {
	defer close(f.out)
	defer close(f.errors)

	for url := range f.in {
		urls, err := f.Crawler.Crawl(url)
		if err != nil {
			f.errors <- err
		}

		f.out <- FetchResult{
			src:     url,
			results: urls,
		}
	}
}

func (f Fetcher) Results() <-chan FetchResult {
	return f.out
}

func (f Fetcher) Errors() <-chan error {
	return f.errors
}
