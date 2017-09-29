package scrape

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type attrs = map[string]string

func TestLinkGathering(t *testing.T) {
	bytes, err := ioutil.ReadFile("testdata/simple-test.html")
	if err != nil {
		t.Fatalf("Unable to read file; %v", err)
	}

	scraper, err := Create(bytes)
	if err != nil {
		t.Fatalf("Unable to create scraper; %v", err)
	}

	body := scraper.FindFirst("p", attrs{"class": "content"})
	links := body.FindAll("a", attrs{"class": "target"})

	if len(links) != 3 {
		t.Fatalf("Expected to match 3 links. Instead found %d", len(links))
	}

	for i, link := range links {
		expected := fmt.Sprintf("http://testlink-%d", i+1)
		actual, ok := link.GetAttr("href")
		if !ok {
			t.Fatalf("Returned tag has no href attribute")
		}
		if actual != expected {
			t.Fatalf("Expected link to equal '%s' but found '%s'", expected, actual)
		}
	}
}
