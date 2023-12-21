package store

import (
	"context"
	"sync"
)

// Compile time check that (*mem) satisfies Store interface
var _ Store = &mem{}

type mem struct {
	sync.RWMutex
	users map[string]struct{}
	items map[string]map[string]string
}

func NewInMemory() Store {
	return &mem{
		users: make(map[string]struct{}),
		items: make(map[string]map[string]string),
	}
}

func (m *mem) UserExists(ctx context.Context, username string) (bool, error) {
	m.RLock()
	defer m.RUnlock()
	_, ok := m.users[username]
	return ok, nil
}

func (m *mem) UserAdd(ctx context.Context, username string) error {
	exists, err := m.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if exists {
		return ErrorUserAlreadyExists
	}

	m.Lock()
	defer m.Unlock()
	m.users[username] = struct{}{}
	m.items[username] = make(map[string]string)
	return nil
}

func (m *mem) ItemGet(ctx context.Context, username, key string) (value string, err error) {
	m.RLock()
	defer m.RUnlock()

	exists, err := m.UserExists(ctx, username)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrorUserNotFound
	}

	value, ok := m.items[username][key]
	if !ok {
		return "", ErrorItemNotFound
	}
	return value, nil
}

func (m *mem) ItemSet(ctx context.Context, username, key, value string) error {

	exists, err := m.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrorUserNotFound
	}

	m.Lock()
	defer m.Unlock()
	m.items[username][key] = value
	return nil
}

func (m *mem) ItemDelete(ctx context.Context, username, key string) error {
	exists, err := m.UserExists(ctx, username)
	if err != nil {
		return err
	}
	if !exists {
		return ErrorUserNotFound
	}

	m.Lock()
	m.Unlock()
	delete(m.items[username], key)
	return nil
}
