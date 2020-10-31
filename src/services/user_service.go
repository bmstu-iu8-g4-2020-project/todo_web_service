package services

import (
	"log"

	"todo_web_service/src/models"
)

func (db *DataBase) AddUser(user models.User) error {
	result, err := db.Exec("INSERT INTO tg_user (username, user_id) values ($1, $2)", user.UserName, user.Id)
	if err != nil {
		return err
	}

	log.Println(result.LastInsertId())
	log.Println(result.RowsAffected())
	return nil
}

func (db *DataBase) UserInfo(userId int) (models.User, error) {
	user := models.User{}
	row := db.QueryRow("SELECT * FROM tg_user WHERE user_id= $1", userId)
	err := row.Scan(&user.Id, &user.UserName)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
