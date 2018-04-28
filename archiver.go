package archiver

import (
	"time"
)

type Feed interface {
	Title() string
	Type() string
	Author() string
	Updated() string
	Link() string
	Description() string
	Language() string
	Items() []FeedItem
	Refresh(Fetcher) (newItems bool, err error)
}

type FeedItem interface {
	Url() string
	Content() string
}

type Repository interface {
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Dirty() bool
	Dir() string // TODO
	Clean() error
}

type ArchiveStore interface {
	Get(key string, locationPath string) error
	Put(repository Repository) (key string, err error)
	IsPinned(hash string) (pinned bool, err error)
	// Pin(hash string) (err error)
	// UnPin(hash string) (err error)
}

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type Cache interface {
	LatestValue(key string) ([]byte, error)
	Get(key string, getter func() ([]byte, error), expires time.Duration) ([]byte, error)
}
