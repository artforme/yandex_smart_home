package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	UserID string
	Hash   string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "284411"
	dbname   = "yandex_smart_home"
)

type Storage struct {
	dataBase *sqlx.DB
}

func New() (*Storage, error) {
	const op = "storage.postgres.New"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: fail conection to the db: %w", op, err)
	}

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
		primary_key SERIAL PRIMARY KEY,
		user_id TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id)
	`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		return nil, fmt.Errorf("%s: users execute statement: %w", op, err)
	}

	createYandexTable := `
		CREATE TABLE IF NOT EXISTS yandex_tokens (
		user_id INTEGER PRIMARY KEY REFERENCES users(primary_key),
		token TEXT NOT NULL
	)`
	_, err = db.Exec(createYandexTable)
	if err != nil {
		return nil, fmt.Errorf("%s: users execute statement: %w", op, err)
	}

	createControllersTable := `
		CREATE TABLE IF NOT EXISTS controller_tokens (
		user_id INTEGER PRIMARY KEY REFERENCES users(primary_key),
		token TEXT NOT NULL
	)`
	_, err = db.Exec(createControllersTable)
	if err != nil {
		return nil, fmt.Errorf("%s: users execute statement: %w", op, err)
	}

	return &Storage{dataBase: db}, nil
}

func (s *Storage) Insert() error {
	return nil
}

func (s *Storage) Search(userID string) error {
	const op = "storage.postgres.Search"

	searchIDUser := `
		 SELECT user_id, password 
		 FROM users 
		 WHERE user_id = $1
	`
	row := s.dataBase.QueryRow(searchIDUser, userID)

	var currentUser User
	err := row.Scan(&currentUser.UserID, &currentUser.Hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return fmt.Errorf("%s: scan result: %w", op, err)
	}

	return nil
}
func (s *Storage) Close() error {
	if err := s.dataBase.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}
	return nil
}
