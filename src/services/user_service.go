package services

import (
	"log"

	"todo_web_service/src/models"
)

type Datastore interface {
	AddUserToDB(user models.User) error
	UserInfo(user models.User) (models.User, error)
}

func (db *DataBase) AddUserToDB(user models.User) error {
	result, err := db.Exec("INSERT INTO tg_user (username, user_id) values ($1, $2)", user.UserName, user.Id)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db *DataBase) UserInfo(user models.User) (models.User, error) {
	row := db.QueryRow("SELECT * FROM tg_user WHERE username= $1", user.UserName)
	err := row.Scan(&user.Id, &user.UserName)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
