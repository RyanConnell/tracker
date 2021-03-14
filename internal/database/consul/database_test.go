package consul

import (
	"errors"
	"strings"

	"github.com/hashicorp/consul/api"
)

var errNotFound = errors.New("kv: not found")

type testKV struct {
	m map[string][]byte
}

func (kv *testKV) Get(key string, _ *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
	if v, exists := kv.m[key]; exists {
		return &api.KVPair{Key: key, Value: v}, nil, nil
	} else {
		return nil, nil, errNotFound
	}
}

func (kv *testKV) Put(pair *api.KVPair, _ *api.WriteOptions) (*api.WriteMeta, error) {
	kv.m[pair.Key] = pair.Value
	return nil, nil
}

func (kv *testKV) List(prefix string, _ *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
	var pairs []*api.KVPair
	for key, value := range kv.m {
		if strings.HasPrefix(key, prefix) {
			pairs = append(pairs, &api.KVPair{Key: key, Value: value})
		}
	}

	return api.KVPairs(pairs), nil, nil
}
