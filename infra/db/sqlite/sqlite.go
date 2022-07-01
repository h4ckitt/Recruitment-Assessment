package sqlite

import (
	"assessment/apperror"
	"assessment/config"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Repo struct {
	db *sql.DB
}

func NewSqliteClient() (*Repo, error) {
	conf := config.FetchConfig()

	connectionInfo := fmt.Sprintf(fmt.Sprintf("%s", conf.DatabaseFileName))

	log.Println(connectionInfo)

	db, err := sql.Open("sqlite3", connectionInfo)

	if err != nil {
		return nil, err
	}

	return &Repo{db}, nil
}

func (repo *Repo) FetchPaginatedPhoneNumbers(offset, limit int) ([]string, error) {
	var (
		result []string
	)

	query := fmt.Sprintf("SELECT phone FROM customer LIMIT %d, %d", offset, limit)

	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var phone string

		if err := rows.Scan(&phone); err != nil {
			return nil, err
		}

		result = append(result, phone)
	}

	return result, nil
}

func (repo *Repo) FetchPaginatedPhoneNumbersByCode(code string, offset, limit int) ([]string, error) {
	var (
		result []string
	)

	query := fmt.Sprintf("SELECT phone FROM customer WHERE phone LIKE '(%s)%%' LIMIT %d, %d", code, offset, limit)

	rows, err := repo.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var phone string

		if err := rows.Scan(&phone); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = apperror.NotFound
			}
			return nil, err
		}

		result = append(result, phone)
	}
	return result, nil
}
