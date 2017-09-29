package show

import (
	"fmt"
)

// Show struct must implement Trackable
type Show struct {
	Name       string
	Episodes   []*Episode
	Finished   bool
	EpisodeURL string
}

type Episode struct {
	Title       string
	Season      string
	Episode     string
	ReleaseDate string
}

func (s *Show) Write() error {
	return nil
}

func (s *Show) String() string {
	episodeString := ""
	for _, episode := range s.Episodes {
		episodeString += "\t" + episode.String() + "\n"
	}

	return fmt.Sprintf("%s: Episodes=%d, EpisodeURL=%s\n%s", s.Name,
		len(s.Episodes), s.EpisodeURL, episodeString)

}

func (s *Episode) String() string {
	return fmt.Sprintf("%sx%s: %s - '%s'", s.Season, s.Episode, s.ReleaseDate, s.Title)
}
