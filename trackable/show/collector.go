package show

import (
	"fmt"
	"strconv"
	"strings"

	"tracker/date"
	"tracker/scrape"
	"tracker/trackable/common"
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
	body, err := common.GetBytes(url)
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
	body, err := common.GetBytes(url)
	if err != nil {
		return err
	}

	scraper, err := scrape.Create(body)
	if err != nil {
		return fmt.Errorf("Unable to create scraper; %v - %v\n", err, scraper)
	}

	seasonNum := 1
	var previousDate *date.Date = &date.Date{}
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

func (s *Show) parseEpisodeTable(table *scrape.Tag, season int, previousDate *date.Date) error {
	rows := table.FindAll("tr", nil)
	for _, row := range rows {
		if !row.Valid {
			continue
		}

		columns := row.FindAll("td", nil)
		if len(columns) < 2 {
			continue
		}

		episodeNumStr := common.ParseString(columns[0].Text())
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
				episode.Title = common.ParseString(column.Text())
				continue
			}

			// Get release date
			text := common.ParseString(column.Text())
			if date.IsDate(text) {
				episode.ReleaseDate, err = date.ToDate(text)
				if err != nil {
					fmt.Printf("Unable to convert %s to a date object: %v\n", text, err)
				}
				continue
			}

		}

		// If new date is less than previous date we can skip this episode.
		// (Indicative of Webisodes)
		if episode.ReleaseDate.CompareTo(previousDate) == -1 {
			continue
		}

		previousDate = episode.ReleaseDate
		s.Episodes = append(s.Episodes, episode)
	}
	return nil
}
