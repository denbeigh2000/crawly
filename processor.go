package crawly

type Processor struct {
	in          <-chan FetchResult
	newURLs     chan string
	crawledURLs chan string

	// Keeps a set of URLs that we've seen before
	seen map[string]struct{}

	// Keeps a set of URLs that are currently in the pipeline
	tracker map[string]struct{}
}

func (p *Processor) loop() {
	defer close(p.newURLs)
	defer close(p.crawledURLs)

	for result := range p.in {
		addedMore := false
		delete(p.tracker, result.src)

		// If the source is coming through this end of the pipeline, it would
		// have had to have come from the other end, so we don't send this
		// down the New URLs channel. This is primarily to address the edge case
		// where the first URL of the entire crawl is fetched twice.
		p.seen[result.src] = struct{}{}

		for _, newUrl := range result.results {
			if _, ok := p.seen[newUrl]; !ok {
				p.seen[newUrl] = struct{}{}
				p.tracker[newUrl] = struct{}{}
				p.newURLs <- newUrl
				addedMore = true
			}
		}

		p.crawledURLs <- result.src

		if !addedMore && len(p.tracker) == 0 {
			// There were no new results from this result, and there are no more URLs
			// in the pipeline, end
			return
		}
	}
}

func NewProcessor(in <-chan FetchResult) *Processor {
	p := &Processor{
		in:          in,
		newURLs:     make(chan string),
		crawledURLs: make(chan string),

		seen:    make(map[string]struct{}),
		tracker: make(map[string]struct{}),
	}

	go p.loop()
	return p
}

func (p *Processor) New() <-chan string {
	return p.newURLs
}

func (p *Processor) Crawled() <-chan string {
	return p.crawledURLs
}
