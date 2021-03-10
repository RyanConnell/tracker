package sql

import (
	"context"
	"database/sql"
	"fmt"

	"tracker/internal/types/show"
)

type ShowsDatabase struct {
	db *Database

	listStmt     *sql.Stmt
	detailsStmt  *sql.Stmt
	episodesStmt *sql.Stmt
}

func (db *Database) Shows() *ShowsDatabase {
	listStmt, err := db.db.Prepare(listQuery)
	if err != nil {
		// TODO: Is there a cleaner way of doing this? Is this OK?
		panic(fmt.Sprintf("unable to prepare list statement: %v", err))
	}

	detailsStmt, err := db.db.Prepare(detailsQuery)
	if err != nil {
		// TODO: Is there a cleaner way of doing this? Is this OK?
		panic(fmt.Sprintf("unable to prepare details statement: %v", err))
	}

	episodesStmt, err := db.db.Prepare(episodesQuery)
	if err != nil {
		// TODO: Is there a cleaner way of doing this? Is this OK?
		panic(fmt.Sprintf("unable to prepare episodes statement: %v", err))
	}

	return &ShowsDatabase{
		db:           db,
		listStmt:     listStmt,
		detailsStmt:  detailsStmt,
		episodesStmt: episodesStmt,
	}
}

func (db *ShowsDatabase) List(ctx context.Context) ([]*show.Show, error) {
	rows, err := db.listStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list shows: %w", err)
	}

	var shows []*show.Show
	for rows.Next() {
		show := &show.Show{}
		if err := rows.Scan(
			&show.ID,
			&show.Name,
			&show.WikipediaURL,
			&show.TrailerURL,
			&show.Finished,
		); err != nil {
			return nil, fmt.Errorf("unable to get show details: %w", err)
		}
		shows = append(shows, show)
	}

	return shows, nil
}

func (db *ShowsDatabase) Details(ctx context.Context, id string) (*show.Show, error) {
	s := &show.Show{}
	if err := db.detailsStmt.QueryRowContext(ctx, id).Scan(
		&s.ID,
		&s.Name,
		&s.WikipediaURL,
		&s.TrailerURL,
		&s.Finished,
	); err != nil {
		return nil, fmt.Errorf("unable to get show details: %w", err)
	}

	rows, err := db.episodesStmt.QueryContext(ctx, s.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to show episodes: %w", err)
	}

	for rows.Next() {
		episode := &show.Episode{}
		if err := rows.Scan(
			&episode.Title,
			&episode.Season,
			&episode.Episode,
			&episode.ReleaseDate,
		); err != nil {
			return nil, fmt.Errorf("unable to get episode: %w", err)
		}
		s.Episodes = append(s.Episodes, episode)
	}

	return s, nil
}

var listQuery = `
SELECT
	id,
	title,
	wikipedia,
	trailer,
	finished,
FROM shows
`

var detailsQuery = listQuery + " WHERE id=? LIMIT 1"

var episodesQuery = `
SELECT
	title,
	season,
	episode,
	release_date
FROM episodes
WHERE show_id=?
`
