package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"yandex_smart_house/internal/random"
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

func (s *Storage) Insert(userID string) error {
	const op = "storage.postgres.Insert"

	insertStr := `
		INSERT INTO users (user_id, password) 
		VALUES ($1, $2);
	`
	userPassword := random.NewRandomString(32)
	fmt.Println(userPassword)
	hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), 10)
	fmt.Println(hash)
	if err != nil {
		return fmt.Errorf("%s: hash generation fail: %w", op, err)
	}

	_, err = s.dataBase.Exec(insertStr, userID, string(hash))
	if err != nil {
		return fmt.Errorf("%s: hash generation fail: %w", op, err)
	}
	return nil
}

func (s *Storage) Search(userID string, userPassword string) error {
	const op = "storage.postgres.Search"

	searchIDUser := `
		SELECT password
		FROM users
		WHERE user_id = $1;
	`

	var hashFromDB []byte

	err := s.dataBase.QueryRow(searchIDUser, userID).Scan(&hashFromDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
		return fmt.Errorf("%s: QueryRow fail: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(hashFromDB, []byte(userPassword))
	if err != nil {
		return fmt.Errorf("%s: password does't match: %w", op, err)
	}

	return nil
}
func (s *Storage) Close() error {
	if err := s.dataBase.Close(); err != nil {
		return fmt.Errorf("error closing database connection: %w", err)
	}
	return nil
}
