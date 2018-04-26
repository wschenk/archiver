package archiver

// import (
// 	"io"
// )

type Feed interface {
	Title() string
	Type() string
	Author() string
	Updated() string
	Link() string
	Description() string
	Language() string
	Items() []FeedItem
	Refresh() (newItems bool, err error)
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
}

type ArchiveStore interface {
	Get(key string, locationPath string) error
	Put(repository Repository) (key string, err error)
	IsPinned(hash string) (pinned bool, err error)
	// Pin(hash string) (err error)
	// UnPin(hash string) (err error)
}
