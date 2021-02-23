package show

import (
	"fmt"

	"tracker/date"
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

type ShowSimple struct {
	ID    int
	Name  string
	Image string
}

type ShowList struct {
	Count int
	Shows []*ShowSimple
}

type ShowFull struct {
	*Show

	SeasonCount  int `json:"season_count"`
	EpisodeCount int `json:"episode_count"`

	// Episode info
	MostRecentEpisode *Episode `json:"most_recent_episode"`
	NextEpisode       *Episode `json:"next_episode"`
}

type Schedule struct {
	StartDate date.Date      `json:"start_date"`
	EndDate   date.Date      `json:"end_date"`
	Items     []ScheduleItem `json:"items"`
}

type ScheduleItem struct {
	Date     *date.Date       `json:"date"`
	Episodes []*CalendarEntry `json:"episodes"`
}

var listFilters = map[string]func(*Show) bool{
	"all":        listFilterAll,
	"airing":     listFilterAiring,
	"upcoming":   listFilterUpcoming,
	"unreleased": listFilterUnreleased,
}

func (h *Handler) Get(id int) (*ShowFull, error) {
	if id <= 0 || len(h.shows) < id {
		return nil, fmt.Errorf("Invalid show ID")
	}
	show := h.shows[id-1]
	sf := &ShowFull{show, 0, 0, show.GetMostRecentEpisode(), show.GetNextEpisode()}
	sf.SeasonCount = sf.MostRecentEpisode.Season
	sf.EpisodeCount = show.EpisodesBefore(sf.MostRecentEpisode) + 1
	return sf, nil
}

func (h *Handler) GetList(listType string) (*ShowList, error) {
	filter, ok := listFilters[listType]
	if !ok {
		return nil, fmt.Errorf("Unknown list type: %s", listType)
	}
	shows := h.shows
	showsSimple := make([]*ShowSimple, 0)
	for _, show := range shows {
		if filter(show) {
			showsSimple = append(showsSimple, showToSimple(show))
		}
	}
	return &ShowList{
		Count: len(showsSimple),
		Shows: showsSimple,
	}, nil
}

func (h *Handler) GetSchedule(start, end string) (*Schedule, error) {
	startDate, err := date.DateFromStr(start)
	if err != nil {
		return nil, err
	}

	endDate, err := date.DateFromStr(end)
	if err != nil {
		return nil, err
	}

	schedule := &Schedule{}
	schedule.StartDate = startDate
	schedule.EndDate = endDate
	dateRange, err := date.DatesInRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Date range: %v", err)
	}

	days := make([]ScheduleItem, len(dateRange))
	episodeMap := h.episodesInRange(dateRange)
	for i, date := range dateRange {
		item := ScheduleItem{
			Date:     date,
			Episodes: episodeMap[*date],
		}
		days[i] = item
	}

	schedule.Items = days
	return schedule, nil
}

type CalendarEntry struct {
	ShowID   int
	ShowName string

	*Episode
}

func (h *Handler) episodesInRange(dateRange []*date.Date) map[date.Date][]*CalendarEntry {
	episodeMap := map[date.Date][]*CalendarEntry{}
	if len(dateRange) == 0 {
		return episodeMap
	}
	for _, date := range dateRange {
		episodeMap[*date] = make([]*CalendarEntry, 0)
	}
	for _, show := range h.shows {
		eps := show.EpisodesInRange(dateRange[0], dateRange[len(dateRange)-1])
		for _, e := range eps {
			entry := &CalendarEntry{show.ID, show.Name, e}
			episodeMap[*e.ReleaseDate] = append(episodeMap[*e.ReleaseDate], entry)
		}
	}
	return episodeMap
}

func showToSimple(show *Show) *ShowSimple {
	s := ShowSimple{
		ID:    show.ID,
		Name:  show.Name,
		Image: show.Image,
	}
	return &s
}

// listFilterAll will always return true.
func listFilterAll(show *Show) bool {
	return true
}

// listFilterAiring will return true for all shows that have aired in the last 30 days.
func listFilterAiring(show *Show) bool {
	lastMonth := date.CurrentDate().Minus(30)
	if mostRecentEpisode := show.GetMostRecentEpisode(); mostRecentEpisode != nil {
		return mostRecentEpisode.ReleaseDate.CompareTo(lastMonth) == 1
	}
	return false
}

// listFilterUpcoming will return true for all upcoming shows.
func listFilterUpcoming(show *Show) bool {
	return !listFilterAiring(show) && show.GetNextEpisode() != nil
}

// listFilterUnreleased will return true for all unreleased shows.
func listFilterUnreleased(show *Show) bool {
	return show.GetMostRecentEpisode() == nil
}
