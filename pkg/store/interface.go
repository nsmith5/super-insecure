package store

import "context"

type Store interface {
	// User operations
	UserExists(ctx context.Context, username string) (exists bool, err error)
	UserAdd(ctx context.Context, username string) error

	// Item operations
	ItemGet(ctx context.Context, username, key string) (value string, err error)
	ItemSet(ctx context.Context, username, key, value string) error
	ItemDelete(ctx context.Context, username, key string) error
}
