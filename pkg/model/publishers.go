package model

import (
	"context"
	"database/sql"
	// "errors"
	"log"
	"time"
)

type Publisher struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Headquarters string `json:"headquarters"`
	Website      string `json:"website"`
}

type PublisherModel struct {
	DB *sql.DB
	InfoLog *log.Logger
	ErrorLog *log.Logger
}


func (m PublisherModel) Get(id int) (*Publisher, error) {
	query := `
		SELECT * 
		FROM publishers
		WHERE id = $1 
	`
	var pub Publisher
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&pub.Id, &pub.Name, &pub.Headquarters, &pub.Website)
	if err != nil {
		return nil, err
	}

	return &pub, nil 
}

func (m PublisherModel) GetAll() ([]Publisher, error) {
	query := `
		SELECT * 
		FROM publishers
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pubs []Publisher

	for rows.Next() {
		var pub Publisher
		err := rows.Scan(&pub.Id, &pub.Name, &pub.Headquarters, &pub.Website)
		if err != nil {
			return pubs, err
		}
		pubs = append(pubs, pub)
	}
	if err := rows.Err(); err != nil {
		return pubs, err
	}
	return pubs, nil
}
