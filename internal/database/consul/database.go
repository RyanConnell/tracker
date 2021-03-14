// Package consul allows to write the user data to consul.
package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/hashicorp/consul/api"

	"tracker/internal/database"
)

// UserDatabase implements database.UserDatabase and can be used
// to insert user data.
type Database struct {
	kv KV

	prefix string
}

func NewDatabase(prefix string, opts ...Option) (*Database, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("unable to create client: %w", err)
	}

	db := &Database{
		kv:     client.KV(),
		prefix: prefix,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}

// get value from the database
func (db *Database) get(ctx context.Context, key string, value interface{}) error {
	opt := &api.QueryOptions{}
	opt = opt.WithContext(ctx)

	p := path.Join(db.prefix, key)

	kv, _, err := db.kv.Get(p, opt)
	if err != nil {
		return fmt.Errorf("unable to fetch %s: %w", p, err)
	}
	if kv == nil {
		return fmt.Errorf("%s not found: %w", p, database.ErrNotFound)
	}

	if err := json.Unmarshal(kv.Value, value); err != nil {
		return fmt.Errorf("unable to unmarshal data: %w", err)
	}

	return nil
}

func (db *Database) put(ctx context.Context, key string, value interface{}) error {
	opt := &api.WriteOptions{}
	opt = opt.WithContext(ctx)

	v, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("unable to encode value: %w", err)
	}

	p := path.Join(db.prefix, key)

	if _, err := db.kv.Put(&api.KVPair{Key: p, Value: v}, opt); err != nil {
		return fmt.Errorf("unable to put data: %w", err)
	}

	return nil
}

// Option allows to set options for the Consul database.type Option func(*Database)
type Option func(*Database)

func KVClient(kv KV) Option {
	return func(db *Database) {
		db.kv = kv
	}
}

// KV abstracts over consul KV interface, for testability
type KV interface {
	Put(p *api.KVPair, q *api.WriteOptions) (*api.WriteMeta, error)
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
}
