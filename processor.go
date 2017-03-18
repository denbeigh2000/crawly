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
			addedMore := false

			delete(p.tracker, result.src)
			p.seen[result.src] = struct{}{}

			for _, newUrl := range result.results {
				if _, ok := p.seen[newUrl]; !ok {
					p.seen[newUrl] = struct{}{}
					p.tracker[newUrl] = struct{}{}
					p.newURLs <- newUrl
					addedMore = true
				}
			}

			// fmt.Printf("Crawled: %v,\n\tprovided: %v,\n\tadded = %v\n", result.src, result.results, addedMore)
			// fmt.Printf("\tTracker: %+v\n", p.tracker)
			// fmt.Printf("\tSeen: %+v\n\n", p.seen)

			p.crawledURLs <- result.src

			if !addedMore && len(p.tracker) == 0 {
				// There were no new results from this result, and there are no more URLs
				// in the pipeline, end
				return
			}
		}
	}()
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
