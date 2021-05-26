package database

import (
	"database/sql"
	"fmt"
	"github.com/lgrossi/go-scaffold/src/logger"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string
	Password string
}

type User struct {
	ID       int64  `id:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type ApiUser struct {
	Email string
	Name  string
}

func CreateUser(db *sql.DB, user *User) *User {
	statement := fmt.Sprintf(
		"INSERT INTO users(email, password, name) VALUES ('%s', '%s', '%s')",
		user.Email,
		hashAndSalt([]byte(user.Password)),
		user.Name,
	)

	res, err := db.Exec(statement)
	if err != nil {
		logger.Error(err)
		return nil
	}

	id, _ := res.LastInsertId()

	return GetUserById(db, id)
}

func Login(db *sql.DB, request *AuthRequest) *User {
	user := User{Email: request.Email}
	statement := fmt.Sprintf(
		"SELECT id, name, password FROM users WHERE email = '%s'",
		user.Email,
	)

	err := db.QueryRow(statement).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil || !comparePasswords([]byte(user.Password), []byte(request.Password)) {
		return nil
	}

	user.Password = ""

	return &user
}

func GetUserById(db *sql.DB, userID int64) *User {
	statement := fmt.Sprintf("SELECT id, email, name FROM users WHERE id = '%d'", userID)

	user := User{ID: userID}
	if err := db.QueryRow(statement).Scan(&user.ID, &user.Email, &user.Name); err != nil {
		return nil
	}

	return &user
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)

	if err != nil {
		logger.Panic(err)
	}

	return string(hash)
}

func comparePasswords(hashedPwd, plainPwd []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPwd, plainPwd) == nil
}
