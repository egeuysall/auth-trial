package supabase

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var Db *pgxpool.Pool

func GetUser(ctx context.Context, user User) (User, error) {
	var foundUser User

	username := strings.ToLower(user.Username)
	email := strings.ToLower(user.Email)

	query := `
		SELECT id, username, email, password
		FROM users
		WHERE LOWER(username) = $1 OR LOWER(email) = $2
	`
	err := Db.QueryRow(ctx, query, username, email).
		Scan(&foundUser.Id, &foundUser.Username, &foundUser.Email, &foundUser.Password)

	return foundUser, err
}

func CreateUser(ctx context.Context, user User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := Db.Exec(ctx, query, strings.ToLower(user.Username), strings.ToLower(user.Email), user.Password)

	return err
}
