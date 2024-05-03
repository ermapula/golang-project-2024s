package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/ermapula/golang-project/pkg/validator"
	"github.com/lib/pq"
)

type Game struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"-"`
	Genres      []string  `json:"genres"`
	ReleaseDate time.Time `json:"releaseDate"`
	Price       float64   `json:"price"`
	PublisherId int       `json:"publisherId"`
	Version     int32     `json:"version"`
}

type GameModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// var games = []Game{
// 	{
// 		Id:          "1",
// 		Title:       "Battlefield V",
// 		Genre:       "First Person Shooter",
// 		ReleaseDate: "2018-11-20",
// 		Price:       "$59.99",
// 		PublisherId: "1",
// 	},
// 	{
// 		Id:          "2",
// 		Title:       "Assassin's Creed Unity",
// 		Genre:       "Action-Adventure",
// 		ReleaseDate: "2014-11-11",
// 		Price:       "$39.99",
// 		PublisherId: "2",
// 	},
// 	{
// 		Id:          "3",
// 		Title:       "The Legend of Zelda: Breath of the Wild",
// 		Genre:       "Action-Adventure",
// 		ReleaseDate: "2017-03-03",
// 		Price:       "$59.99",
// 		PublisherId: "3",
// 	},
// 	{
// 		Id:          "4",
// 		Title:       "Call of Duty: Warzone",
// 		Genre:       "Battle Royale",
// 		ReleaseDate: "2020-03-10",
// 		Price:       "0",
// 		PublisherId: "4",
// 	},
// 	{
// 		Id:          "5",
// 		Title:       "Elden Ring",
// 		Genre:       "Action RPG",
// 		ReleaseDate: "2022-02-25",
// 		Price:       "$59.99",
// 		PublisherId: "5",
// 	},
// 	{
// 		Id:          "6",
// 		Title:       "Apex Legends",
// 		Genre:       "Battle Royale",
// 		ReleaseDate: "2019-02-04",
// 		Price:       "0",
// 		PublisherId: "1",
// 	},
// 	{
// 		Id:          "7",
// 		Title:       "Far Cry 6",
// 		Genre:       "First Person Shooter",
// 		ReleaseDate: "2021-10-07",
// 		Price:       "$59.99",
// 		PublisherId: "2",
// 	},
// 	{
// 		Id:          "8",
// 		Title:       "Super Mario Odyssey",
// 		Genre:       "Platformer",
// 		ReleaseDate: "2017-10-27",
// 		Price:       "$49.99",
// 		PublisherId: "3",
// 	},
// 	{
// 		Id:          "9",
// 		Title:       "Call of Duty: Ghosts",
// 		Genre:       "First Person Shooter",
// 		ReleaseDate: "2014-03-25",
// 		Price:       "$59.99",
// 		PublisherId: "4",
// 	},
// 	{
// 		Id:          "10",
// 		Title:       "Dark Souls III",
// 		Genre:       "Action RPG",
// 		ReleaseDate: "2019-03-22",
// 		Price:       "$59.99",
// 		PublisherId: "5",
// 	},
// }

func (m GameModel) Get(id int) (*Game, error) {
	query := `
		SELECT id, created_at, title, genres, price, release_date, publisher_id, version
		FROM games
		WHERE id = $1 
	`
	var game Game
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&game.Id, 
		&game.CreatedAt,
		&game.Title, 
		pq.Array(&game.Genres), 
		&game.Price, 
		&game.ReleaseDate, 
		&game.PublisherId,
		&game.Version,
	)
	if err != nil {
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &game, nil
}

func (m GameModel) Post(game *Game) error {
	query := `
		INSERT INTO games (title, genres, price, release_date, publisher_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version
	`
	args := []interface{}{game.Title, pq.Array(game.Genres), game.Price, game.ReleaseDate, game.PublisherId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&game.Id, &game.CreatedAt, &game.Version); err != nil {
		return err
	}

	return nil
}

func (m GameModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM games WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)

	return err
}

func (m GameModel) Update(game *Game) error {
	query := `
		UPDATE games
		SET title = $1, genres = $2, price = $3, release_date = $4, publisher_id = $5, version = version + 1
		WHERE id = $6
		RETURNING version
	`

	args := []interface{}{
		game.Title, 
		pq.Array(game.Genres), 
		game.Price,
		game.ReleaseDate, 
		game.PublisherId, 
		game.Id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.Version)
}

func ValidateGame(v *validator.Validator, game *Game) {
	v.Check(game.Title != "", "title", "must be provided")
	v.Check(len(game.Title) <= 500, "title", "must not be more than 500 bytes long")
}
