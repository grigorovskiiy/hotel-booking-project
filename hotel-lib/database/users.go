package database

import (
	"fmt"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/models"
	"github.com/google/uuid"
	"log"
)

func Exists(user models.UserLogin, withPassword bool) (bool, error) {
	query := NewSqlBuilder()

	var condition string
	if withPassword {
		condition = fmt.Sprintf("login = '%s' and password = '%s'", user.Login, user.Password)
	} else {
		condition = fmt.Sprintf("login = '%s'", user.Login)
	}
	query = query.Select([]string{"COUNT(*)"}).
		From("users").
		Where(condition)

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return false, err
	}
	defer db.Close()

	log.Print(query)
	result, err := db.Exists(query)
	if err != nil {
		return false, err
	}
	return result, nil
}

func AddUser(user models.User) error {
	query := NewSqlBuilder()
	query = query.CustomQuery(fmt.Sprintf("INSERT INTO users (id, role, login, password) VALUES ('%s', '%s', '%s', '%s')",
		uuid.NewString(), user.Role, user.Login, user.Password)).Returning("id")

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return err
	}
	defer db.Close()

	log.Print(query)
	_, err = db.Insert(query)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(login string) (models.User, error) {
	query := NewSqlBuilder()
	query = query.Select([]string{}).From("users").Where(fmt.Sprintf("login = '%s'", login))

	db, err := InitConnection("hotel-bookings")
	if err != nil {
		return models.User{}, err
	}
	defer db.Close()

	log.Print(query)
	var user models.User

	rows, err := db.GetAll(query)
	if err != nil {
		return models.User{}, err
	}
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Role, &user.Login, &user.Password)
		if err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}
