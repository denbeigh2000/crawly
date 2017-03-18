# crawly

This is a simple web crawler exercise written in Go.

The "Crawler" implementation, which fetches reachable URLs given a starting URL is currently
statically defined as a map, but I intend to make a simple HTML-parser driven one by fetching
the URL and parsing the resultant HTML for \<a\> tags
