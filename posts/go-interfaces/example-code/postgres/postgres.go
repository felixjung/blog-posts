package postgres

import (
	"context"
	"fmt"

	"github.com/felixjung/blog-post/go-interfaces/example-code/server"
	"github.com/jackc/pgx/v4"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(config *pgx.ConnConfig) (*DB, error) {
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return &DB{conn}, nil
}

func (db *DB) Create(u server.User) error {
	// TODO: Create the user in postgres
	return nil
}

func (db *DB) Read(userID string) (server.User, error) {
	// TODO: Read the user from postgres
	return server.User{}, nil
}

func (db *DB) Update(u server.User) (server.User, error) {
	// TODO: Update the user in postgres.
	return server.User{}, nil
}

func (db *DB) Delete(userID string) (server.User, error) {
	// TODO: Delete the user in postgres.
	return server.User{}, nil
}
