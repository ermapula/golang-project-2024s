package model

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"fmt"

	"github.com/ermapula/golang-project/pkg/validator"
	"github.com/lib/pq"
)

type Game struct {
	Id          int64    `json:"id"`
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

func (m GameModel) GetAll(title string, genres []string, publisher_id int, filters Filters) ([]*Game, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, genres, price, release_date, publisher_id, version
		FROM games
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		AND (publisher_id = $3 OR $3 = -1)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
	`, filters.sortColumn(), filters.sortDirection())
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	args := []interface{}{title, pq.Array(genres), publisher_id, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var games []*Game

	for rows.Next() {
		var game Game
		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}
		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return games, metadata, nil
}

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

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.Id, &game.CreatedAt, &game.Version)
}


func (m GameModel) Update(game *Game) error {
	query := `
		UPDATE games
		SET title = $1, genres = $2, price = $3, release_date = $4, publisher_id = $5, version = version + 1
		WHERE id = $6 AND version = $7
		RETURNING version
	`

	args := []interface{}{
		game.Title, 
		pq.Array(game.Genres), 
		game.Price,
		game.ReleaseDate, 
		game.PublisherId, 
		game.Id,
		game.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&game.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
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

func ValidateGame(v *validator.Validator, game *Game) {
	v.Check(game.Title != "", "title", "must be provided")
	v.Check(len(game.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(game.ReleaseDate.Before(time.Now()), "releaseDate", "must be a date before today")
	v.Check(game.Price >= 0, "price", "must be at least zero")
	v.Check(game.PublisherId > 0, "publisherId", "must be a positive integer")
	v.Check(len(game.Genres) > 0, "genres", "must contain at least one genre")
}

func (m GameModel) GetAllOfUser(userId int64) ([]*Game, error) {
	query := `
		SELECT g.id, g.title, g.created_at, g.genres, g.price, g.release_date, g.publisher_id, g.version
		FROM games g
		JOIN library l ON g.id = l.game_id
		WHERE l.user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []*Game
	for rows.Next() {
		var game Game
		err := rows.Scan(
			&game.Id, 
			&game.Title, 
			&game.CreatedAt, 
			pq.Array(&game.Genres), 
			&game.Price, 
			&game.ReleaseDate, 
			&game.PublisherId,
			&game.Version,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, &game)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

func (m GameModel) AddToLibrary(userId int64, gameId int) error {
	query := `
		SELECT id
		FROM library
		WHERE user_id = $1 AND game_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var id int 
	row := m.DB.QueryRowContext(ctx, query, userId, gameId)
	err := row.Scan(&id)
	if err == nil {
		return errors.New("game already in library")
	}
	if err != sql.ErrNoRows {
		return err
	}


	query = `
		INSERT INTO library (user_id, game_id)
		VALUES ($1, $2)
	`
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.DB.ExecContext(ctx, query, userId, gameId)
	
	return err
}

func (m GameModel) DeleteFromLibrary(userId int64, gameId int) error {
	if userId < 1 || gameId < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM library WHERE user_id = $1 AND game_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, userId, gameId)

	return err
}