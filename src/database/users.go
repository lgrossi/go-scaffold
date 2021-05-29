package database

import (
	"crypto/hmac"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/lgrossi/go-scaffold/src/logger"
	"golang.org/x/crypto/sha3"
)

const (
	secret = "MY_FUCKING_GREAT_SECRET_YOLO"
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
	hexPass, _ := hex.DecodeString(user.Password)

	if err != nil || !comparePasswords(hexPass, []byte(request.Password)) {
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
	h := hmac.New(sha3.New512, []byte(secret))
	h.Write(pwd)

	return hex.EncodeToString(h.Sum(nil))
}

func comparePasswords(hashedPwd, plainPwd []byte) bool {
	h := hmac.New(sha3.New512, []byte(secret))
	h.Write(plainPwd)
	return hmac.Equal(h.Sum(nil), hashedPwd)
}
