package services

import (
	"database/sql"
	"fmt"
)

type DataBase struct {
	*sql.DB
}

func NewDB(dbUser string, dbPassword string) (*DataBase, error) {
	dbSourceName := fmt.Sprintf("user=%s password=%s dbname=todownik sslmode=disable", dbUser, dbPassword)

	db, err := sql.Open("postgres", dbSourceName)
	if err != nil {
		return nil, err
	}
	// check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database!")
	return &DataBase{db}, nil
}
