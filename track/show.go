package track

import ()

// Show struct must implement Trackable
type Show struct {
	Name     string
	Episodes []Episode
	Finished bool
}

type Episode struct {
	Title string
}

func (s *Show) Write() error {
	return nil
}
