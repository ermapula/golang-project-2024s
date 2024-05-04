package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Publishers PublisherModel
	Games      GameModel
	Users      UserModel
	Tokens     TokenModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	return Models{
		Publishers: PublisherModel{
			DB:       db,
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		},
		Games: GameModel{
			DB:       db,
			ErrorLog: errorLog,
			InfoLog:  infoLog,
		},
		Users: UserModel{
			DB: db,
		},
		Tokens: TokenModel{
			DB: db,
		},
	}
}
