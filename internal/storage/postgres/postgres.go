package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/KKitsun/link-shortener/internal/config"
)

type Storage struct {
	db *sql.DB
}

func NewPostgres(cfg *config.Config) (*Storage, error) {

	dblInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.DBname)
	db, err := sql.Open("postgres", dblInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS link(
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON link(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("NewPostgres.Exec: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {

	stmt, err := s.db.Prepare("INSERT INTO link(url, alias) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("SaveURL.stmt prepearing error: %w", err)
	}

	var insertedRowId int64

	row := stmt.QueryRow(urlToSave, alias)
	if err := row.Scan(&insertedRowId); err != nil {
		return 0, fmt.Errorf("SaveURL. stmt execution insert error: %w", err)
	}

	return insertedRowId, nil
}

func (s *Storage) GetURL(alias string) (string, error) {

	stmt, err := s.db.Prepare("SELECT url FROM link WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("GetURL. stmt prepearing error: %w", err)
	}

	var recievedURL string

	err = stmt.QueryRow(alias).Scan(&recievedURL)
	if err != nil {
		return "", fmt.Errorf("GetURL. stmt execution select error: %w", err)
	}

	return recievedURL, nil
}
