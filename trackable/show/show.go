package show

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Show struct must implement Trackable
type Show struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Episodes     []*Episode `json:"-"`
	Finished     bool       `json:"finished"`
	EpisodeURL   string     `json:"episode_url"`
	WikipediaURL string     `json:"wikipedia"`
	TrailerURL   string     `json:"trailer"`

	// Backwards Compatability
	Location string `json:"location"`
	Airing   int    `json:"airing"`
	Upcoming int    `json:"upcoming"`
	Image    string `json:"image"`
	ImdbUrl  string `json:"imdb_url"`
}

type Episode struct {
	Title       string
	Season      int
	Episode     int
	ReleaseDate string
}

func (s *Show) Write() error {
	db, err := openDB("tracker")
	if err != nil {
		return err
	}

	for _, e := range s.Episodes {
		_, err = db.Exec(`INSERT INTO episodes(show_id, season, episode, title, release_date)
		 		          VALUES(?, ?, ?, ?, ?)`, s.ID, e.Season, e.Episode, e.Title, e.ReleaseDate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Show) String() string {
	episodeString := ""
	for _, episode := range s.Episodes {
		episodeString += "\t" + episode.String() + "\n"
	}

	return fmt.Sprintf("%-2d - %-30s - %3d Episodes, WikipediaURL='%s'\n%s", s.ID, s.Name,
		len(s.Episodes), s.WikipediaURL, episodeString)

}

func (s *Episode) String() string {
	return fmt.Sprintf("%10s x %-50s: %-50s - '%-50s'", s.Season, s.Episode,
		s.ReleaseDate, s.Title)
}

func ScrapeAll() error {
	shows, err := loadAllShows()
	if err != nil {
		return err
	}

	errors := make([]error, 0)

	for _, show := range shows {
		fmt.Printf("\t%v", show)
		url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", show.WikipediaURL)
		if err = show.scrape(url); err != nil {
			return err
		}
		err = show.Write()
		if err != nil {
			err = fmt.Errorf("Error with show %d; %v", show.ID, err)
			errors = append(errors, err)
		}
	}

	fmt.Printf("%d Error(s) occured\n", len(errors))
	for _, err := range errors {
		fmt.Printf("Error: %v\n", err)
	}

	return nil
}

func (s *Show) Scan(rows *sql.Rows) error {
	return rows.Scan(&s.ID, &s.Name, &s.WikipediaURL, &s.TrailerURL,
		&s.Finished)
}

func loadAllShows() ([]*Show, error) {
	shows := make([]*Show, 0)

	db, err := openDB("tracker")
	if err != nil {
		return shows, err
	}

	rows, err := db.Query("SELECT id,title,wikipedia,trailer,finished FROM shows")
	if err != nil {
		return shows, err
	}

	for rows.Next() {
		show := &Show{}
		err := show.Scan(rows)
		if err != nil {
			return shows, err
		}
		shows = append(shows, show)
	}

	return shows, nil
}

func openDB(name string) (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("rhino:@/%s", name))
	if err != nil {
		return nil, err
	}
	return db, nil
}
