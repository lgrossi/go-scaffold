package database

import (
	"database/sql"
	"fmt"
)

type User struct {
	ID       int64
	Email    string
	Password string
}

func Login(db *sql.DB, user *User) bool {
	statement := fmt.Sprintf(
		"SELECT id FROM users WHERE email = '%s' and password = '%s'",
		user.Email,
		user.Password,
	)

	err := db.QueryRow(statement).Scan(&user.ID)
	return err == nil
}

func GetUserById(db *sql.DB, userID int64) *User {
	statement := fmt.Sprintf("SELECT email, password FROM users WHERE id = '%d'", userID)

	user := User{ID: userID}
	if err := db.QueryRow(statement).Scan(&user.Email, &user.Password); err != nil {
		return nil
	}

	return &user
}
