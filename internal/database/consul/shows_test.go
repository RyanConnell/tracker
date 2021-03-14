package consul

import (
	"context"
	"sort"
	"testing"
	"tracker/internal/types/show"

	"github.com/go-test/deep"
)

var testData = map[string][]byte{
	"tracker/shows/list/westworld":            []byte(`{"id": 1, "name": "Westworld"}`),
	"tracker/shows/episodes/westworld_s01e01": []byte(`{"title":"The Original", "season": 1, "episode": 1}`),

	"tracker/shows/list/expanse": []byte(`{"id": 2, "name": "The Expanse"}`),
}

func TestList(t *testing.T) {
	kv := &testKV{m: testData}
	db, err := NewDatabase("tracker", KVClient(kv))
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}

	showsDB := db.Shows()

	want := []*show.Show{
		{ID: 2, Name: "The Expanse"},
		{ID: 1, Name: "Westworld"},
	}

	got, err := showsDB.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error listing shows: %v", err)
	}

	sort.Sort(show.ByName(got))

	if diff := deep.Equal(got, want); diff != nil {
		t.Fatalf("List() = %v, want %v, diff = %v", got, want, diff)
	}
}

func TestDetails(t *testing.T) {
	kv := &testKV{m: testData}
	db, err := NewDatabase("tracker", KVClient(kv))
	if err != nil {
		t.Fatalf("unable to setup database: %v", err)
	}

	showsDB := db.Shows()

	want := &show.Show{
		ID:   1,
		Name: "Westworld",
		Episodes: []*show.Episode{
			{Title: "The Original", Season: 1, Episode: 1},
		},
	}

	got, err := showsDB.Details(context.Background(), "westworld")
	if err != nil {
		t.Fatalf("unexpected error getting details: %v", err)
	}

	if diff := deep.Equal(got, want); diff != nil {
		t.Fatalf("Details() = %v, want %v, diff = %v", got, want, diff)
	}

}
