package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

// Compile time check that (*mem) satisfies Store interface
var _ Store = &db{}

type db struct {
	*pgx.Conn
}

func NewPostgres(conn *pgx.Conn) Store {
	return &db{conn}
}

func (d *db) UserExists(ctx context.Context, username string) (bool, error) {
	query := `
	SELECT count(id)
	FROM users
	WHERE username = $1
	`
	var count int
	err := d.QueryRow(ctx, query, username).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 1 {
		return false, errors.New("data corruption")
	}
	return count == 1, nil
}

func (d *db) UserAdd(ctx context.Context, username string) error {
	query := `
	INSERT into users (username) VALUES ($1)
	`
	tag, err := d.Exec(ctx, query, username)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return ErrorUserAlreadyExists
	}
	return nil
}

func (d *db) ItemGet(ctx context.Context, username, key string) (value string, err error) {
	query := `
	SELECT (value)
	FROM items
	INNER JOIN users ON items.user_id = users.id
	WHERE users.username = $1 AND items.name = $2
	`
	err = d.QueryRow(ctx, query, username, key).Scan(&value)
	return
}

func (d *db) ItemSet(ctx context.Context, username, key, value string) error {
	query := `SELECT id FROM users WHERE username = $1`
	var userId int
	err := d.QueryRow(ctx, query, username).Scan(&userId)
	if err != nil {
		return err
	}

	query = `INSERT into items (user_id, name, value) VALUES ($1, $2, $3)`
	tag, err := d.Exec(ctx, query, userId, key, value)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("should have changed one row")
	}
	return nil
}

func (d *db) ItemDelete(ctx context.Context, username, key string) error {
	query := `SELECT id FROM users WHERE username = $1`
	var userId int
	err := d.QueryRow(ctx, query, username).Scan(&userId)
	if err != nil {
		return err
	}

	query = `DELETE FROM items
	WHERE user_id = $1 AND name = $2`
	tag, err := d.Exec(ctx, query, userId, key)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("should have changed one row")
	}
	return nil
}
