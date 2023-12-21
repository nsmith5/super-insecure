package store

import (
	"context"
	"testing"
)

type predicate func(ctx context.Context, s Store) (bool, error)

func predicateUserExists(username string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		return s.UserExists(ctx, username)
	}
}

func predicateUserAddFails(username string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		err := s.UserAdd(ctx, username)
		return err != nil, nil
	}
}

func predicateItemGetFails(username, key string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		_, err := s.ItemGet(ctx, username, key)
		return err != nil, nil
	}
}

func predicateItemSetFails(username, key, value string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		err := s.ItemSet(ctx, username, key, value)
		return err != nil, nil
	}
}

func predicateItemDeleteFails(username, key string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		err := s.ItemDelete(ctx, username, key)
		return err != nil, nil
	}
}

func predicateItemIs(username, key, value string) predicate {
	return func(ctx context.Context, s Store) (bool, error) {
		got, err := s.ItemGet(ctx, username, key)
		if err != nil {
			return false, err
		}
		return got == value, nil
	}
}

func TestInMemory(t *testing.T) {
	tests := map[string]struct {
		state      func(ctx context.Context, s Store) error
		predicates map[string]predicate
	}{
		"can't add the same user twice": {
			state: func(ctx context.Context, s Store) error {
				return s.UserAdd(ctx, "alice")
			},
			predicates: map[string]predicate{
				"alice exists":                   predicateUserExists("alice"),
				"adding alice again should fail": predicateUserAddFails("alice"),
			},
		},
		"can't do item stuff for users that don't exist": {
			state: func(ctx context.Context, s Store) error {
				return nil
			},
			predicates: map[string]predicate{
				"can't set item":    predicateItemSetFails("alice", "foo", "bar"),
				"can't get item":    predicateItemGetFails("alice", "foo"),
				"can't delete item": predicateItemDeleteFails("alice", "foo"),
			},
		},
		"item get after set works": {
			state: func(ctx context.Context, s Store) error {
				err := s.UserAdd(ctx, "alice")
				if err != nil {
					return err
				}
				err = s.ItemSet(ctx, "alice", "foo", "bar")
				if err != nil {
					return err
				}
				return nil
			},
			predicates: map[string]predicate{
				"alice's foo is bar": predicateItemIs("alice", "foo", "bar"),
			},
		},
		"different users have different items": {
			state: func(ctx context.Context, s Store) error {
				err := s.UserAdd(ctx, "alice")
				if err != nil {
					return err
				}
				err = s.ItemSet(ctx, "alice", "foo", "bar")
				if err != nil {
					return err
				}

				err = s.UserAdd(ctx, "bob")
				if err != nil {
					return err
				}
				err = s.ItemSet(ctx, "bob", "foo", "baz")
				if err != nil {
					return err
				}

				return nil
			},
			predicates: map[string]predicate{
				"alice's foo is bar": predicateItemIs("alice", "foo", "bar"),
				"bob's foo is baz":   predicateItemIs("bob", "foo", "baz"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := NewInMemory()
			ctx := context.Background()

			err := test.state(ctx, s)
			if err != nil {
				t.Fatal(err)
			}

			for pname, predicate := range test.predicates {
				passed, err := predicate(ctx, s)
				if err != nil {
					t.Fatal(err)
				}
				if !passed {
					t.Errorf("predicate %q failed", pname)
				}
			}
		})
	}
}
