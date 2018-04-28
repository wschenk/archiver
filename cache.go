package archiver

import (
	"fmt"
	"time"
)

type RepoCache struct {
	repo Repository
}

func CreateRepoCache(repo Repository) *RepoCache {
	return &RepoCache{repo}
}

func (cache *RepoCache) LatestValue(key string) ([]byte, error) {
	return cache.repo.Get(key)
}

func (cache *RepoCache) Get(key string, getter func() ([]byte, error), expires time.Duration) ([]byte, error) {
	data, err := cache.LatestValue(key)

	if err != nil {
		return nil, err
	}

	if data != nil {
		fmt.Printf("Return %s from cache\n", key)
		return data, nil
	}

	fmt.Printf("Loading from function\n")
	data, err = getter()

	if err != nil {
		return nil, err
	}

	if data != nil {
		err := cache.repo.Put(key, data)

		if err != nil {
			return data, err
		}
	}

	return data, nil
}
