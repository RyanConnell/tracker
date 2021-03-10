package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"tracker/internal/types/show"
)

type ShowsDatabase struct {
	db     *Database
	prefix string
}

func (db *Database) Shows() *ShowsDatabase {
	return &ShowsDatabase{
		db:     db,
		prefix: "shows",
	}
}

func (db *ShowsDatabase) List(ctx context.Context) ([]*show.Show, error) {
	showIDs, err := db.list(ctx, "list")
	if err != nil {
		return nil, fmt.Errorf("unable to get show IDs")
	}

	// TODO: Speed up by running these in parallel
	var shows []*show.Show
	for showID := range showIDs {
		show, err := db.showDetails(ctx, showID)
		if err != nil {
			return nil, fmt.Errorf("unable to get show %s: %w", showID, err)
		}
		shows = append(shows, show)
	}

	return shows, nil
}

func (db *ShowsDatabase) Details(ctx context.Context, id string) (*show.Show, error) {
	// TODO: Run getting show details and show episodes in parallel
	s, err := db.showDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("unable to get details about %s: %w", id, err)
	}

	episodes, err := db.list(ctx, path.Join("episodes", id))
	if err != nil {
		return nil, fmt.Errorf("unable to get episodes for %s: %w", id, err)
	}

	for _, ep := range episodes {
		var episode show.Episode
		if err := json.Unmarshal(ep, &episode); err != nil {
			return nil, fmt.Errorf("unable to parse episode: %w", err)
		}
		s.Episodes = append(s.Episodes, &episode)
	}

	return s, nil
}

func (db *ShowsDatabase) showDetails(ctx context.Context, id string) (*show.Show, error) {
	var show show.Show
	if err := db.get(ctx, path.Join("list", id), &show); err != nil {
		return nil, err
	}

	return &show, nil
}

func (db *ShowsDatabase) get(ctx context.Context, key string, value any) error {
	return db.db.get(ctx, path.Join(db.prefix, key), value)
}

func (db *ShowsDatabase) list(ctx context.Context, prefix string) (map[string][]byte, error) {
	return db.db.list(ctx, path.Join(db.prefix, prefix))
}
