package services

import (
	"database/sql"
	"fmt"
)

type DataBase struct {
	*sql.DB
}

func NewDB(dbSourceName string) (*DataBase, error) {
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
