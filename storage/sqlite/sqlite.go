package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"main.go/lib/e"
	"main.go/storage"
)

type DataStorage struct {
	db *sql.DB
}

func New(pathToDb string) (*DataStorage, error) {
	db, err := sql.Open("sqlite3", pathToDb)

	if err != nil {
		return nil, fmt.Errorf("can't open Database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't establish connection with Database: %w", err)
	}

	return &DataStorage{db: db}, nil
}

func (ds *DataStorage) Save(ctx context.Context, userData *storage.Userdata) error {
	q := `INSERT INTO pages (user_name, city) VALUES (?, ?)`

	if _, err := ds.db.ExecContext(ctx, q, userData.UserName, userData.City); err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}

	return nil
}

func (ds *DataStorage) Remove(ctx context.Context, userData storage.Userdata) error {
	q := `DELETE FROM pages WHERE user_name = ?`

	if _, err := ds.db.ExecContext(ctx, q, userData.UserName); err != nil {
		return fmt.Errorf("can't delete data: %w", err)
	}

	return nil
}

func (ds *DataStorage) Exists(ctx context.Context, userData storage.Userdata) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE user_name = ?`

	var count int

	if err := ds.db.QueryRowContext(ctx, q, userData.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists %w", err)
	}

	return count > 0, nil
}

func (ds *DataStorage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (user_name TEXT, city TEXT)`
	_, err := ds.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (ds *DataStorage) RetrieveCity(ctx context.Context, userData storage.Userdata) (string, error) {
	exists, err := ds.Exists(ctx, userData)
	if err != nil {
		return "", fmt.Errorf("can't check if page exists %w", err)
	}

	if !exists {
		return "", nil
	}

	q := `SELECT city FROM pages WHERE user_name = ?`

	var result string

	err = ds.db.QueryRowContext(ctx, q, userData.UserName).Scan(&result)
	if err != nil {
		return "", e.Wrap("can't prepare sql query", err)
	}

	return result, nil
}
