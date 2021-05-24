package database

import (
	"database/sql"
	"fmt"
)

type User struct {
	id       int64
	Email    string
	Password string
	// other fields here
}

func Login(db *sql.DB, user *User) bool {
	statement := fmt.Sprintf(
		"SELECT id FROM users WHERE email = '%s' and password = '%s'",
		user.Email,
		user.Password,
	)

	err := db.QueryRow(statement).Scan(&user.id)
	return err == nil
}
