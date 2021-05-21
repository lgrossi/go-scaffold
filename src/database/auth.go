package database

import (
	"database/sql"
	"fmt"
	"github.com/lgrossi/go-scaffold/src/logger"
	"time"
)

type ActiveToken struct {
	Email     string
	TokenStr  string
	ExpiresAt int64
}

func IsActiveToken(db *sql.DB, token *ActiveToken) bool {
	statement := fmt.Sprintf(
		"SELECT expires_at FROM active_tokens WHERE email = '%s' and token = '%s'",
		token.Email,
		token.TokenStr,
	)

	if err := db.QueryRow(statement).Scan(&token.ExpiresAt); err != nil {
		return false
	}

	return token.ExpiresAt <= time.Now().Unix()
}

func InsertActiveToken(db *sql.DB, token ActiveToken) bool {
	statement := fmt.Sprintf(
		`INSERT INTO active_tokens(email, token, expires_at) VALUES ("%s", "%s", %d)`,
		token.Email,
		token.TokenStr,
		token.ExpiresAt,
	)

	if _, err := db.Exec(statement); err != nil {
		logger.Error(err)
		return false
	}

	return true
}

func DeactivateToken(db *sql.DB, token ActiveToken) bool {
	statement := fmt.Sprintf(
		`DELETE FROM active_tokens WHERE token = "%s" AND email = "%s"`,
		token.TokenStr,
		token.Email,
	)

	if _, err := db.Exec(statement); err != nil {
		logger.Error(err)
		return false
	}

	return true
}
