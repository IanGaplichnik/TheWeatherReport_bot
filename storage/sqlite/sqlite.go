package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"main.go/clients/events"
	"main.go/lib/e"
	"main.go/storage"
)

type DataStorage struct {
	db *sql.DB
}

const (
	MenuNot      = "none"
	MenuKeyboard = "keyboard"
)

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
	exists, err := ds.Exists(ctx, userData)
	if err != nil {
		return e.Wrap("can't save", err)
	}

	var query string

	if !exists {
		query = `INSERT INTO pages (user_name) VALUES (?)`
		if _, err := ds.db.ExecContext(ctx, query, userData.UserName); err != nil {
			return fmt.Errorf("can't save: %w", err)
		}
	}

	query = `UPDATE pages SET city=?, menu=?, country=?, state=? WHERE user_name=?`
	if _, err := ds.db.ExecContext(ctx,
		query,
		userData.City,
		userData.Menu,
		userData.Country,
		userData.State,
		userData.UserName); err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}

	return nil
}

func (ds *DataStorage) Remove(ctx context.Context, userData *storage.Userdata) error {
	q := `DELETE FROM pages WHERE user_name = ?`

	if _, err := ds.db.ExecContext(ctx, q, userData.UserName); err != nil {
		return fmt.Errorf("can't delete data: %w", err)
	}

	return nil
}

func (ds *DataStorage) Exists(ctx context.Context, userData *storage.Userdata) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE user_name = ?`

	var count int

	if err := ds.db.QueryRowContext(ctx, q, userData.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists %w", err)
	}

	return count > 0, nil
}

func (ds *DataStorage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (user_name TEXT, city TEXT, menu TEXT, country TEXT, state TEXT)`
	_, err := ds.db.ExecContext(ctx, q)

	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}

func (ds *DataStorage) RetrieveLocation(ctx context.Context, userData *storage.Userdata) (*events.CityData, error) {
	exists, err := ds.Exists(ctx, userData)
	if err != nil {
		return nil, fmt.Errorf("can't check if page exists %w", err)
	}

	if !exists {
		return nil, nil
	}

	var result events.CityData

	q := `SELECT city, country, state FROM pages WHERE user_name = ?`
	err = ds.db.QueryRowContext(ctx, q, userData.UserName).Scan(&result.CityName, &result.Country, &result.State)
	if err != nil {
		return nil, e.Wrap("can't prepare sql query", err)
	}

	return &result, nil
}

func (ds *DataStorage) IsKeyboardMenu(ctx context.Context, userData *storage.Userdata) (bool, error) {
	query := `SELECT menu FROM pages WHERE user_name=?`

	var result string

	err := ds.db.QueryRowContext(ctx, query, userData.UserName).Scan(&result)
	if err != nil {
		return false, e.Wrap("can't prepare sql query", err)
	}

	if result == MenuNot {
		return false, nil
	}
	return true, nil
}
