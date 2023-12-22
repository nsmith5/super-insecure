package store

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestPostgres(t *testing.T) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("skipping postgres tests. Set DATABASE_URL to run them")
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close(ctx)

	s := NewPostgres(conn)

	_ = s.UserAdd(ctx, "alice")
	exists, err := s.UserExists(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Error("where is alice?")
	}

	err = s.ItemSet(ctx, "alice", "foo", "bar")
	if err != nil {
		t.Error(err)
	}

	value, err := s.ItemGet(ctx, "alice", "foo")
	if err != nil {
		t.Error(err)
	}
	if value != "bar" {
		t.Error("what?")
	}

	err = s.ItemDelete(ctx, "alice", "foo")
	if err != nil {
		t.Error(err)
	}
}
