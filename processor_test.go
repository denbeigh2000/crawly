package crawly

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessorBasic(t *testing.T) {
	assert := assert.New(t)

	in := make(chan FetchResult)
	p := NewProcessor(in)

	src := "www.somecompany.com"
	dest := "www.somecompany.com/about"

	entry := FetchResult{
		src:     src,
		results: []string{dest},
	}

	in <- entry

	newURLs := p.New()
	crawledURLs := p.Crawled()

	newEntry := <-newURLs
	srcEntry := <-crawledURLs

	assert.Equal(newEntry, dest, "Newly-found URL should be sent down the channel")
	assert.Equal(srcEntry, src, "Source URL should be sent down the crawled channel")

	close(in)

	_, ok := <-newURLs
	assert.False(ok, "Channel should be closed after processing finished")

	_, ok = <-crawledURLs
	assert.False(ok, "Channel should be closed after processing finished")
}

func TestProcessorMultipleURLs(t *testing.T) {
	assert := assert.New(t)

	src := "www.somecompany.com"
	dests := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
	}

	entry := FetchResult{
		src:     src,
		results: dests,
	}

	in := make(chan FetchResult)
	p := NewProcessor(in)

	in <- entry

	crawled := p.Crawled()
	new := p.New()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for _, dest := range dests {
			newDest := <-new
			assert.Equal(newDest, dest, "URLs should be sent in order")
		}

		wg.Done()
	}()

	go func() {
		newSrc := <-crawled
		assert.Equal(newSrc, src, "Crawled channel should send source URL")
		wg.Done()
	}()

	close(in)

	wg.Wait()
	_, ok := <-new
	assert.False(ok, "Channel should be closed after processing finished")

	_, ok = <-crawled
	assert.False(ok, "Channel should be closed after processing finished")
}

func TestProcessorRepeatedURLs(t *testing.T) {
	assert := assert.New(t)

	src := "www.somecompany.com"
	dests := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
	}

	expectedDests := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
	}

	entry := FetchResult{
		src:     src,
		results: dests,
	}

	in := make(chan FetchResult)
	p := NewProcessor(in)

	in <- entry

	crawled := p.Crawled()
	new := p.New()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for _, dest := range expectedDests {
			newDest := <-new
			assert.Equal(newDest, dest, "URLs should be sent in order")
		}

		wg.Done()
	}()

	go func() {
		newSrc := <-crawled
		assert.Equal(newSrc, src, "Crawled channel should send source URL")
		wg.Done()
	}()

	close(in)

	wg.Wait()
	_, ok := <-new
	assert.False(ok, "Channel should be closed after processing finished")

	_, ok = <-crawled
	assert.False(ok, "Channel should be closed after processing finished")
}

func TestProcessorRepeatedURLsMultipleMessages(t *testing.T) {
	assert := assert.New(t)

	src1 := "www.somecompany.com"
	dests1 := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
	}

	src2 := "www.somecompany.com/jobs"
	dests2 := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/only-on-the-second-link",
	}

	expectedDests := []string{
		"www.somecompany.com/about",
		"www.somecompany.com/jobs",
		"www.somecompany.com/dmca",
		"www.somecompany.com/only-on-the-second-link",
	}

	in := make(chan FetchResult)
	p := NewProcessor(in)

	go func() {
		in <- FetchResult{
			src:     src1,
			results: dests1,
		}

		in <- FetchResult{
			src:     src2,
			results: dests2,
		}

		close(in)
	}()

	crawled := p.Crawled()
	new := p.New()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for _, dest := range expectedDests {
			newDest := <-new
			assert.Equal(newDest, dest, "URLs should be sent in order")
		}

		wg.Done()
	}()

	go func() {
		assert.Equal(<-crawled, src1, "Crawled channel should send source URL")
		assert.Equal(<-crawled, src2, "Crawled channel should send source URL")
		wg.Done()
	}()

	wg.Wait()
	_, ok := <-new
	assert.False(ok, "Channel should be closed after processing finished")

	_, ok = <-crawled
	assert.False(ok, "Channel should be closed after processing finished")
}
