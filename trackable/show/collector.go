package show

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tracker/internal/timeutil"
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
	body, err := getBytes(url)
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

	if s.EpisodeURL == "" {
		s.EpisodeURL = url
	}
	if err := s.scrapeEpisodes(s.EpisodeURL); err != nil {
		return err
	}

	fmt.Printf("Show: %+v\n", s)
	return nil
}

func (s *Show) scrapeEpisodes(url string) error {
	body, err := getBytes(url)
	if err != nil {
		return err
	}

	scraper, err := scrape.Create(body)
	if err != nil {
		return fmt.Errorf("Unable to create scraper; %v - %v\n", err, scraper)
	}

	seasonNum := 1
	var previousDate time.Time
	tables := scraper.FindAll("table", attr{"class": "wikiepisodetable"})
	for _, table := range tables {
		if !table.Valid {
			continue
		}

		err := s.parseEpisodeTable(table, seasonNum, previousDate)
		seasonNum++
		if err != nil {
			return err
		}

	}
	return nil
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

func (s *Show) parseEpisodeTable(table *scrape.Tag, season int, previousDate time.Time) error {
	rows := table.FindAll("tr", nil)
	for _, row := range rows {
		if !row.Valid {
			continue
		}

		columns := row.FindAll("td", nil)
		if len(columns) < 2 {
			continue
		}

		episodeNumStr := parseString(columns[0].Text())
		episodeNum, err := strconv.Atoi(episodeNumStr)
		if err != nil {
			return fmt.Errorf("Unable to convert %s to an integer: %v", episodeNumStr, err)
		}

		episode := &Episode{
			Season:  season,
			Episode: episodeNum,
		}

		for _, column := range columns {
			if !column.Valid {
				continue
			}

			// Get title
			class, ok := column.GetAttr("class")
			if ok && strings.Contains(class, "summary") {
				episode.Title = parseString(column.Text())
				continue
			}

			// Get release date
			text := parseString(column.Text())
			if timeutil.HasMonth(text) {
				episode.ReleaseDate, err = timeutil.Parse(text)
				if err != nil {
					fmt.Printf("Unable to convert %s to a date object: %v\n", text, err)
				}
				continue
			}

		}

		// If new date is less than previous date we can skip this episode.
		// (Indicative of Webisodes)
		if episode.ReleaseDate.Before(previousDate) {
			continue
		}

		previousDate = episode.ReleaseDate
		s.Episodes = append(s.Episodes, episode)
	}
	return nil
}

func getBytes(url string) ([]byte, error) {
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

func parseString(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\t")
	str = strings.Trim(str, " ")
	str = strings.Replace(str, "  ", " ", 100)
	str = strings.Replace(str, "\t\t", "\t", 100)
	str = strings.Replace(str, "\n", "", 100)
	return str
}
