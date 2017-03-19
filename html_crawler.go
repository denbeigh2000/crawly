package crawly

import (
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var ValidLink = regexp.MustCompile(`^http(s)?://.*`)
var DefaultClient = &http.Client{Timeout: 3 * time.Second}

type HTMLCrawler struct{}

func (h HTMLCrawler) Crawl(url string) ([]string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	linkNodes := document.Find("a").Nodes
	return extractLinks(linkNodes), nil
}

func validLink(link string) bool {
	return ValidLink.MatchString(link)
}

func extractLinks(nodes []*html.Node) []string {
	out := make([]string, len(nodes))
NODE_LOOP:
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" && validLink(attr.Val) {
				out = append(out, attr.Val)
				continue NODE_LOOP
			}
		}
	}

	return out
}
