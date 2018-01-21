package show

// Handler will take care of database loading and API prepping for Shows.
type Handler struct {
	shows []*Show
}

func (h *Handler) Init() {
	shows, err := loadAllShows()
	if err != nil {
		h.shows = make([]*Show, 0)
	} else {
		h.shows = shows
	}
}

func (h *Handler) Get(id int) (*Show, error) {
	return h.shows[id-1], nil
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
