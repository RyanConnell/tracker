package scrape

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type Tag struct {
	token html.Token
	bytes []byte
	Valid bool
}

func Create(bytes []byte) (*Tag, error) {
	scraper := &Tag{
		token: html.Token{},
		bytes: bytes,
	}
	return scraper, nil
}

// FindFirst will return the first matching Tag
func (s *Tag) FindFirst(tag string, params map[string]string) *Tag {
	tags := findTags(s.bytes, tag, params, 1)
	if len(tags) == 0 {
		return &Tag{Valid: false}
	}
	return tags[0]
}

// FindAll will return all matching Tags
func (s *Tag) FindAll(tag string, params map[string]string) []*Tag {
	return findTags(s.bytes, tag, params, -1)
}

// findTags will return "count" matching Tags
func findTags(bytes []byte, tag string, params map[string]string, count int) []*Tag {
	tags := make([]*Tag, 0)

	tokenizer := html.NewTokenizer(strings.NewReader(string(bytes)))

	for len(tags) < count || count == -1 {
		tagType := tokenizer.Next()

		if tagType == html.ErrorToken || len(tags) == count {
			return tags
		}

		if tagType == html.StartTagToken {
			currentTag := tokenizer.Token()
			if currentTag.Data == tag {

				// Return a "tag" object instead of just a html token
				tagData := &Tag{
					token: currentTag,
					bytes: tagContents(currentTag, tokenizer),
					Valid: true,
				}

				// Verify that our required params match the params on the tag.
				if !tagData.paramMatch(params) {
					continue
				}

				tags = append(tags, tagData)
			}
		}
	}

	return tags
}

// Text will retrieve all text from inside a tag
func (t *Tag) Text() string {
	text := ""
	tokenizer := html.NewTokenizer(strings.NewReader(string(t.bytes)))

	for {
		tagType := tokenizer.Next()
		if tagType == html.TextToken {
			text += tokenizer.Token().Data
		} else if tagType == html.ErrorToken {
			break
		}
	}
	return text
}

// tagContents returns the HTML contained within the current Tag
func tagContents(token html.Token, tokenizer *html.Tokenizer) []byte {
	// Start at a given tag and work your way down until the depth gets back to 0.
	bytes := make([]byte, 0)
	depth := 1

	bytes = append(bytes, []byte(fmt.Sprintf("%v", token))...)
	for {
		if depth == 0 {
			break
		}

		tagType := tokenizer.Next()

		if tagType == html.StartTagToken {
			depth++
		} else if tagType == html.EndTagToken {
			depth--
		} else if tagType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		bytes = append(bytes, []byte(fmt.Sprintf("%+v\n", token))...)
		if tagType == html.ErrorToken {
			break
		}

	}
	return bytes
}

var looseMatchParams = map[string]bool{"class": true}

func (t *Tag) paramMatch(params map[string]string) bool {
	for key, value := range params {
		val, ok := t.GetAttr(key)

		// Check for loose matches, such as when an attr can have multiple values.
		if _, useLooseMatch := looseMatchParams[key]; useLooseMatch {
			if looseMatch(val, value) {
				continue
			}
			return false
		}

		// Check for normal matches where attr is exactly equal to the given value.
		if !ok || val != value {
			return false
		}
	}
	return true
}

// looseMatch will return true if {actual} matches one of the items contained in {expected}
// for example: 'test' is matched in 'this is a test' but not in 'thisisatest'
func looseMatch(expected, actual string) bool {
	items := strings.Split(expected, " ")
	for _, item := range items {
		if item == actual {
			return true
		}
	}
	return false
}

// GetAttr will return the value of a specific attribute for the current Tag
func (t *Tag) GetAttr(attr string) (string, bool) {
	for _, attribute := range t.token.Attr {
		if attribute.Key == attr {
			return attribute.Val, true
		}
	}
	return "", false
}

// String returns a string representation of a Tag
func (t *Tag) String() string {
	return fmt.Sprintf("token: %v, bytes: %s", t.token, string(t.bytes))
}
