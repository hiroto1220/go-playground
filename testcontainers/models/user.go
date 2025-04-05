package handler

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    string
	Name  string
	Email string
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FetchUserByID(userID string) (*User, error) {
	row := r.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザーが見つからなかった場合
		}
		return nil, err
	}

	return &user, nil

}
