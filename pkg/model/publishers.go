package model

import (
	"context"
	"database/sql"
	"errors"
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

var publishers = []Publisher{
	{
		Id:           "1",
		Name:         "Electronic Arts",
		Headquarters: "Redwood City, California, USA",
		Website:      "https://www.ea.com",
	},
	{
		Id:           "2",
		Name:         "Ubisoft",
		Headquarters: "Montreuil, France",
		Website:      "https://www.ubisoft.com",
	},
	{
		Id:           "3",
		Name:         "Nintendo",
		Headquarters: "Kyoto, Japan",
		Website:      "https://www.nintendo.com",
	},
	{
		Id:           "4",
		Name:         "Activision Blizzard",
		Headquarters: "Santa Monica, California, USA",
		Website:      "https://www.activisionblizzard.com",
	},
	{
		Id:           "5",
		Name:         "FromSoftware",
		Headquarters: "Tokyo, Japan",
		Website:      "https://www.fromsoftware.jp",
	},
}

// func GetPublishers() []Publisher {
// 	return publishers
// }

// func GetPublisher(id string) (*Publisher, error) {
// 	for _, pub := range publishers {
// 		if pub.Id == id {
// 			return &pub, nil
// 		}
// 	}
// 	return nil, errors.New("Publisher not found")
// }

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