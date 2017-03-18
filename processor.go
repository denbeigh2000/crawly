package crawly

type Processor struct {
	in          <-chan FetchResult
	newURLs     chan string
	crawledURLs chan string

	seen    map[string]struct{}
	tracker map[string]struct{}
}

func (p *Processor) loop() {
	go func() {
		defer close(p.newURLs)
		defer close(p.crawledURLs)

		for result := range p.in {
			for _, newUrl := range result.results {
				delete(p.tracker, result.src)

				addedMore := false
				if _, ok := p.seen[newUrl]; !ok {
					p.newURLs <- newUrl
					p.tracker[newUrl] = struct{}{}
					addedMore = true
				}

				//
				if !addedMore && len(p.tracker) == 0 {
					// There were no new results from this result, and there are no more URLs
					// in the pipeline, end

					return
				}
			}

			p.seen[result.src] = struct{}{}
			p.crawledURLs <- result.src
		}
	}()
}

func NewProcessor(in <-chan FetchResult) *Processor {
	p := &Processor{
		in:          in,
		newURLs:     make(chan string),
		crawledURLs: make(chan string),
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
