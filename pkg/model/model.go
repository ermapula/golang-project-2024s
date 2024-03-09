package model

import (
	"database/sql"
	"log"
	"os"
)

type Models struct {
	Publishers PublisherModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	return Models{
		Publishers: PublisherModel{
			DB: db,
			ErrorLog: errorLog,
			InfoLog: infoLog,
		},
	}
}