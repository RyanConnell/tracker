package track

import (
	"time"
)

// Artist struct must implement Trackable
type Artist struct {
	Name       string
	Albums     []Album
	SoloArtist bool
}

type Album struct {
	Name   string
	Tracks []string
}

type Song struct {
	Title   string
	Runtime time.Duration
}

func (a *Artist) Write() error {
	return nil
}
