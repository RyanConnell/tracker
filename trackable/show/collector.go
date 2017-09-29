package show

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"tracker/scrape"
)

type attr = map[string]string

func Collect() error {
	// This should handle getting a list of all Shows we want and crawling through each one.
	// We shouldn't be giving this a static URL but for now it's fine.
	// In the end it should be able to queue up jobs to be run and the "collector" will
	// take care of them.
	show := &Show{}
	return show.scrape("https://en.wikipedia.org/wiki/Game_of_Thrones")
}

func (s *Show) scrape(url string) error {
	body, err := GetBytes(url)
	if err != nil {
		return err
	}

	scraper, err := scrape.Create(body)
	if err != nil {
		return fmt.Errorf("Unable to create scraper; %v\n", err)
	}

	// Scrape info from the infobox.
	infobox := scraper.FindFirst("table", attr{"class": "infobox"})
	if infobox.Valid {
		if err := s.parseInfobox(infobox); err != nil {
			return err
		}
	}

	if err := s.scrapeEpisodes(s.EpisodeURL); err != nil {
		return err
	}

	fmt.Printf("Show: %+v\n", s)
	return nil
}

// TODO: Make this actually fucking work...
func (s *Show) scrapeEpisodes(url string) error {
	// Scrape info from the Episode List page.
	body, err := GetBytes(url)
	if err != nil {
		return err
	}

	scraper, err := scrape.Create(body)
	if err != nil {
		return fmt.Errorf("Unable to create scraper; %v - %v\n", err, scraper)
	}

	// Scrape info from each episode table
	tables := scraper.FindAll("table", attr{"class": "wikiepisodetable"})
	for _, table := range tables {
		if !table.Valid {
			continue
		}

		rows := table.FindAll("tr", nil)
		for _, row := range rows {
			if !row.Valid {
				continue
			}

			columns := table.FindAll("td", nil)
			if len(columns) < 2 {
				continue
			}

			episode := &Episode{
				Season: "00",
				//Episode: parseString(columns[0].Text()),
			}

			for _, column := range columns {
				if !column.Valid {
					continue
				}

				// Get Title
				class, ok := column.GetAttr("class")
				if ok && strings.Contains(class, "summary") {
					episode.Title = parseString(column.Text())
					if episode.Title[0] == '"' {
						episode.Title = parseString(episode.Title[1 : len(episode.Title)-1])
					}
					continue
				}

				// Get Release Date by checking for a month or datestamp
				text := parseString(column.Text())
				if isDate(text) {
					//episode.ReleaseDate = text
				}

			}

			s.Episodes = append(s.Episodes, episode)

		}

	}

	return nil
}

func parseString(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, " ")
	return str
}

var months = []string{
	"january", "february", "march", "april", "may", "june", "july",
	"august", "september", "october", "november", "december",
}

func isDate(text string) bool {
	for _, month := range months {
		if strings.Contains(strings.ToLower(text), month) {
			return true
		}
	}
	return false
}

func (s *Show) parseInfobox(infobox *scrape.Tag) error {
	if !infobox.Valid {
		return fmt.Errorf("Infobox is not a valid object")
	}

	// Find the name of the Show
	title := infobox.FindFirst("th", attr{"class": "summary"})
	if title.Valid {
		s.Name = strings.Trim(title.Text(), "\n")
	}

	rows := infobox.FindAll("tr", nil)
	for _, row := range rows {
		// Find the link to the episode list
		link := row.FindFirst("a", nil)
		attrib, ok := link.GetAttr("title")
		if ok && strings.Contains(attrib, "List of") {
			if href, ok := link.GetAttr("href"); ok {
				s.EpisodeURL = fmt.Sprintf("http://en.wikipedia.org%s", href)
			}
		}
	}

	return nil
}

func (s *Show) parseEpisodeTable(table *scrape.Tag) error {
	return nil
}

func GetBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Convert reader to bytes
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
