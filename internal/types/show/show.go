package show

import (
	"strings"
	"tracker/trackable/show"
)

// TODO: Move this to a concrete implementation in the future
type Show = show.Show

type Episode = show.Episode

// ByName allows to sort the shows alphabetically by show name.
type ByName []*Show

func (ss ByName) Len() int      { return len(ss) }
func (ss ByName) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
func (ss ByName) Less(i, j int) bool {
	return strings.ToLower(ss[i].Name) < strings.ToLower(ss[j].Name)
}
