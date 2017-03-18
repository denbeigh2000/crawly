package crawly

import (
	"github.com/denbeigh2000/crawly/mock_crawly"

	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBasicFetching(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crawler := mock_crawly.NewMockCrawler(ctrl)

	startPage := "www.somecompany.com"
	targetPages := []string{
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/contact",
	}

	crawler.EXPECT().Crawl(startPage).Return(targetPages, nil).Times(1)

	urls := make(chan string)

	fetcher := NewFetcher(urls, crawler)

	urls <- startPage
	close(urls)

	results := fetcher.Results()

	result := <-results
	assert.Equal(result.src, startPage, "Should persist the source page")
	assert.Equal(result.results, targetPages, "Should send all pages in order")

	_, ok := <-fetcher.Errors()
	assert.False(ok, "Error channel should have closed without error")

	_, ok = <-results
	assert.False(ok, "Result channel should have closed")
}

func TestDuplicatedURLFetching(t *testing.T) {
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crawler := mock_crawly.NewMockCrawler(ctrl)

	startPage := "www.somecompany.com"
	targetPages := []string{
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/about",
		"www.somecompany.com/contact",
		"www.somecompany.com/contact",
		"www.somecompany.com/contact",
		"www.somecompany.com/contact",
		"www.somecompany.com/contact",
	}

	expectedPages := []string{
		"www.somecompany.com/jobs",
		"www.somecompany.com/about",
		"www.somecompany.com/contact",
	}

	crawler.EXPECT().Crawl(startPage).Return(targetPages, nil).Times(1)

	urls := make(chan string)

	fetcher := NewFetcher(urls, crawler)

	urls <- startPage
	close(urls)

	results := fetcher.Results()

	result := <-results
	assert.Equal(result.src, startPage, "Should persist the source page")
	assert.Equal(result.results, expectedPages, "Should send deduplicated pages in order")

	_, ok := <-fetcher.Errors()
	assert.False(ok, "Error channel should have closed without error")

	_, ok = <-results
	assert.False(ok, "Result channel should have closed")
}

func TestErrorReporting(t *testing.T) {
	err := assert.AnError
	assert := assert.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	crawler := mock_crawly.NewMockCrawler(ctrl)

	startPage := "www.somecompany.com"

	crawler.EXPECT().Crawl(startPage).Return(nil, err).Times(1)

	urls := make(chan string)

	fetcher := NewFetcher(urls, crawler)

	urls <- startPage
	close(urls)

	results := fetcher.Results()
	errors := fetcher.Errors()

	er := <-errors
	assert.Equal(err, er, "Should forward the given error")

	result := <-results
	assert.Equal(result.src, startPage, "Should forward the result URL for tracking")
	assert.Equal(result.results, []string(nil), "Should send nil/empty list")

	_, ok := <-errors
	assert.False(ok, "Error channel should have closed")

	_, ok = <-results
	assert.False(ok, "Result channel should have closed")
}
