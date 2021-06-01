package database

import (
	"crypto/hmac"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/lgrossi/go-scaffold/src/logger"
	"golang.org/x/crypto/sha3"
)

const (
	secret = "MY_FUCKING_GREAT_SECRET_YOLO"
)

type LoginRequest struct {
	Email    string
	Password string
}

type RegisterRequest struct {
	Name string
	LoginRequest
}

type User struct {
	ID       int64  `id:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"-"`
	Verified bool   `json:"verified"`
}

type ApiUser struct {
	Email string
	Name  string
}

func CreateUser(db *sql.DB, request *RegisterRequest) *User {
	statement := fmt.Sprintf(
		"INSERT INTO users(email, password, name) VALUES ('%s', '%s', '%s')",
		request.Email,
		hashAndSalt([]byte(request.Password)),
		request.Name,
	)

	res, err := db.Exec(statement)
	if err != nil {
		logger.Error(err)
		return nil
	}

	id, _ := res.LastInsertId()

	return GetUserById(db, id)
}

func Login(db *sql.DB, request *LoginRequest) *User {
	user := GetUserByEmail(db, request.Email)
	if user == nil {
		return nil
	}

	hexPass, _ := hex.DecodeString(user.Password)

	if !comparePasswords(hexPass, []byte(request.Password)) {
		return nil
	}

	return user
}

func SetUserEmailAsVerified(db *sql.DB, email string) bool {

	statement := fmt.Sprintf(
		"UPDATE users SET verified = 1 WHERE email = '%s'", email,
	)

	res, err := db.Exec(statement)
	if err != nil {
		logger.Panic(err)
	}

	rows, _ := res.RowsAffected()

	return rows != 0
}

func ResetPassword(db *sql.DB, email string) string {
	var chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-"
	password := make([]byte, 16)

	_, err := rand.Read(password)
	if err != nil {
		logger.Panic(err)
	}

	for i := 0; i < 16; i++ {
		password[i] = chars[int(password[i])%len(chars)]
	}

	statement := fmt.Sprintf(
		"UPDATE users SET password = '%s' WHERE email = '%s'", hashAndSalt(password), email,
	)

	res, err := db.Exec(statement)
	if err != nil {
		logger.Panic(err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ""
	}

	return string(password)
}

func GetUserById(db *sql.DB, userID int64) *User {
	statement := fmt.Sprintf("SELECT email, name, verified FROM users WHERE id = '%d'", userID)

	user := User{ID: userID}
	if err := db.QueryRow(statement).Scan(&user.Email, &user.Name, &user.Verified); err != nil {
		return nil
	}

	return &user
}

func GetUserByEmail(db *sql.DB, email string) *User {
	user := User{Email: email}
	statement := fmt.Sprintf(
		"SELECT id, name, password, verified FROM users WHERE email = '%s'",
		user.Email,
	)

	err := db.QueryRow(statement).Scan(&user.ID, &user.Name, &user.Password, &user.Verified)

	if err != nil {
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
