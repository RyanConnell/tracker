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
	return h.shows[id], nil
}

func (h *Handler) GetList(count, page int) ([]*Show, error) {
	offset := count * page
	return h.shows[offset : offset+count], nil
}
