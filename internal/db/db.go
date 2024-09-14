package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitDB() {
	connStr := "user=youruser dbname=authdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func StoreRefreshToken(userID, refreshToken string) error {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2)", userID, string(hashedToken))
	return err
}

func ValidateRefreshToken(userID, token string) (bool, error) {
	var storedToken string
	err := db.QueryRow("SELECT token FROM refresh_tokens WHERE user_id = $1", userID).Scan(&storedToken)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedToken), []byte(token))
	return err == nil, err
}
