// Copyright 2020 aaaaaaaalesha <sks2311211@mail.ru>

package services

import (
	"fmt"
	"log"

	"github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/models"
)

func (db *DataBase) AddUser(user models.User) error {
	result, err := db.Exec("INSERT INTO tg_user (user_id, username, state_code, state_request) values ($1, $2, $3, $4);",
		user.Id, user.UserName, user.StateCode, user.StateRequest)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	return nil
}

func (db *DataBase) GetUsers() ([]models.User, error) {
	rows, err := db.Query("SELECT * FROM tg_user;")
	if err != nil {
		return nil, err
	}

	var users []models.User
	for rows.Next() {
		user := models.User{}

		err = rows.Scan(&user.Id, &user.UserName, &user.StateCode, &user.StateRequest)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (db *DataBase) GetUser(userId int) (models.User, error) {
	user := models.User{}
	row := db.QueryRow("SELECT * FROM tg_user WHERE user_id= $1;", userId)
	err := row.Scan(&user.Id, &user.UserName, &user.StateCode, &user.StateRequest)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *DataBase) UpdateState(user models.User) error {
	result, err := db.Exec("UPDATE tg_user SET state_code = $1, state_request = $2 WHERE user_id = $3;",
		user.StateCode, user.StateRequest, user.Id)
	if err != nil {
		return err
	}

	fmt.Println(result.RowsAffected())
	return nil
}
