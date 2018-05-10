package show

import (
	"fmt"
	"strconv"

	"tracker/database"
	"tracker/server/auth"
	"tracker/trackable/common"
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
	ID       int
	Name     string
	Image    string
	UserInfo *UserInfo
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
	StartDate common.Date    `json:"start_date"`
	EndDate   common.Date    `json:"end_date"`
	Items     []ScheduleItem `json:"items"`
}

type ScheduleItem struct {
	Date     *common.Date     `json:"date"`
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

func (h *Handler) GetList(listType string, user *auth.User) (*ShowList, error) {
	filter, ok := listFilters[listType]
	if !ok {
		return nil, fmt.Errorf("Unknown list type: %s", listType)
	}
	shows := h.shows
	showsSimple := make([]*ShowSimple, 0)

	showIDs := make([]int, 0)
	for _, show := range shows {
		showIDs = append(showIDs, show.ID)
	}

	var userInfoList = make(map[int]*UserInfo, 0)
	fmt.Printf("* GetList(%s, %v)\n", listType, user)
	if user != nil {
		GetUserInfo(user.Username, showIDs)
	}

	for _, show := range shows {
		if filter(show) {
			simpleShow := showToSimple(show)
			if userInfo, ok := userInfoList[show.ID]; ok {
				simpleShow.UserInfo = userInfo
			}
			showsSimple = append(showsSimple, simpleShow)
		}
	}
	return &ShowList{
		Count: len(showsSimple),
		Shows: showsSimple,
	}, nil
}

func (h *Handler) GetSchedule(start, end string) (*Schedule, error) {
	startDate, err := common.DateFromStr(start)
	if err != nil {
		return nil, err
	}

	endDate, err := common.DateFromStr(end)
	if err != nil {
		return nil, err
	}

	schedule := &Schedule{}
	schedule.StartDate = startDate
	schedule.EndDate = endDate
	dateRange, err := common.DatesInRange(startDate, endDate)
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

func (h *Handler) episodesInRange(dateRange []*common.Date) map[common.Date][]*CalendarEntry {
	episodeMap := map[common.Date][]*CalendarEntry{}
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

func (h *Handler) RequestShow(u *auth.User, title, wikipedia, trailer, img string) (bool, error) {
	db, err := database.Open("tracker")
	if err != nil {
		return false, err
	}
	defer db.Close()

	if title == "" {
		return false, fmt.Errorf("No show given")
	}

	fmt.Printf("Requested: %s, %s, %s, %s\n", title, wikipedia, trailer, img)
	rows, err := db.Query(`SELECT id FROM requests WHERE title=? OR wikipedia=?`, title, wikipedia)
	if err != nil {
		return false, fmt.Errorf("Error retrieving currently requested shows: %v", err)
	}
	if rows.Next() {
		var id int
		rows.Scan(&id)
		return false, fmt.Errorf("Show already requested: %d", id)
	}

	_, err = db.Exec(`INSERT INTO requests(username, title, wikipedia, trailer, cover_image)
			 VALUES(?, ?, ?, ?, ?)`, u.Username, title, wikipedia, trailer, img)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (h *Handler) TrackShow(u *auth.User, show_id_str, state string) (bool, error) {
	db, err := database.Open("tracker")
	if err != nil {
		return false, err
	}
	defer db.Close()

	if u.Username == "" {
		return false, fmt.Errorf("User must be logged in")
	} else if show_id_str == "" {
		return false, fmt.Errorf("No show id provided")
	} else if state != "true" && state != "false" {
		return false, fmt.Errorf("Expected either 'true' or 'false'")
	}

	show_id, err := strconv.Atoi(show_id_str)
	if err != nil {
		return false, fmt.Errorf("Unable to convert show_id to an integer: %v", err)
	}

	fmt.Printf("TrackShow: %s - %s\n", show_id, state)
	rows, err := db.Query(`SELECT 1 FROM tracked WHERE user_id=? AND show_id=?`,
		u.Username, show_id)
	if err != nil {
		return false, fmt.Errorf("Error retrieving currently tracked shows: %v", err)
	}

	var tracked bool = state == "true"

	// If we already have an entry in the database, update it instead.
	if rows.Next() {
		fmt.Println("Updating entry...")
		_, err := db.Exec("UPDATE tracked SET tracked=? WHERE user_id=? AND show_id=?",
			u.Username, show_id, tracked)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Otherwise insert into the table.
	fmt.Println("Inserting entry...")
	_, err = db.Exec(`INSERT INTO tracked(user_id, show_id, tracked) VALUES(?, ?, ?)`,
		u.Username, show_id, tracked)
	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUserInfo(userID string, showIDs []int) (map[int]*UserInfo, error) {
	fmt.Printf("Attempting to GetUserInfo(%s, %v)...\n", userID, showIDs)
	var userInfoList = make(map[int]*UserInfo, 0)

	db, err := database.Open("tracker")
	if err != nil {
		return userInfoList, err
	}
	defer db.Close()

	// Create a list of ?'s for the SQL search.
	params := []interface{}{userID}
	queryRows := ""
	for i := 0; i < len(showIDs); i++ {
		if i != 0 {
			queryRows += ","
		}
		queryRows += "?"
		params = append(params, showIDs[0])
	}

	// Search for all entries that match the given user_id and any show_id provided.
	query := fmt.Sprintf("SELECT show_id, tracked FROM tracked WHERE user_id=? AND show_id IN (%s)",
		queryRows)
	fmt.Printf("Executing Query: %s\n", query)
	rows, err := db.Query(query, params...)
	if err != nil {
		return userInfoList, err
	}

	for rows.Next() {
		userInfo := &UserInfo{}
		err := userInfo.Scan(rows)
		if err != nil {
			return userInfoList, err
		}
		userInfoList[userInfo.showID] = userInfo
		fmt.Printf("] UserInfo: %v\n", userInfo)
	}

	return userInfoList, nil
}

// listFilterAll will always return true.
func listFilterAll(show *Show) bool {
	return true
}

// listFilterAiring will return true for all shows that have aired in the last 30 days.
func listFilterAiring(show *Show) bool {
	lastMonth := common.CurrentDate().Minus(30)
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
