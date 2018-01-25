package show

import (
	"fmt"
)

// Handler will take care of database loading and API prepping for Shows.
type Handler struct {
	shows []*Show
}

func (h *Handler) Init() {
	shows, err := loadAllShows()
	if err != nil {
		h.shows = make([]*Show, 0)
		fmt.Println(err)
	} else {
		h.shows = shows
	}
}

type ShowFull struct {
	*Show

	SeasonCount  int `json:"season_count"`
	EpisodeCount int `json:"episode_count"`

	// Episode info
	MostRecentEpisode *Episode `json:"most_recent_episode"`
	NextEpisode       *Episode `json:"next_episode"`
}

func (h *Handler) Get(id int) (*ShowFull, error) {
	if id <= 0 || len(h.shows) <= id {
		return nil, fmt.Errorf("Invalid show ID")
	}
	show := h.shows[id-1]
	sf := &ShowFull{show, 0, 0, show.GetMostRecentEpisode(), show.GetNextEpisode()}
	sf.SeasonCount = sf.MostRecentEpisode.Season
	sf.EpisodeCount = show.EpisodesBefore(sf.MostRecentEpisode) + 1
	return sf, nil
}

type ShowSimple struct {
	ID    int
	Name  string
	Image string
}

type ShowList struct {
	Count int
	Shows []*ShowSimple
}

func (h *Handler) GetList() ShowList {
	shows := h.shows
	showsSimple := make([]*ShowSimple, len(shows))
	for i, show := range shows {
		showsSimple[i] = showToSimple(show)
	}
	return ShowList{
		Count: len(showsSimple),
		Shows: showsSimple,
	}
}

func showToSimple(show *Show) *ShowSimple {
	s := ShowSimple{
		ID:    show.ID,
		Name:  show.Name,
		Image: show.Image,
	}
	return &s
}
