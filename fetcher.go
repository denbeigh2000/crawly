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
		out:     make(chan FetchResult, 20),
		errors:  make(chan error, 5),
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

		var results []string
		if err != nil {
			f.errors <- err
			results = nil
		} else {
			results = filterUniqueStrings(urls)
		}

		f.out <- FetchResult{
			src:     url,
			results: results,
		}
	}
}

func (f Fetcher) Results() <-chan FetchResult {
	return f.out
}

func (f Fetcher) Errors() <-chan error {
	return f.errors
}

func filterUniqueStrings(in []string) (out []string) {
	seen := make(map[string]struct{}, len(in))

	for _, url := range in {
		if _, ok := seen[url]; !ok {
			out = append(out, url)
			seen[url] = struct{}{}
		}
	}

	return
}
