package user

import (
	"backend/utils"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	Pool *pgxpool.Pool
}

func checkUserExistsDB(ctx context.Context, tx pgx.Tx, username string) (bool, *utils.HTTPError) {
	var exists bool
	if err := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists); err != nil {
		fmt.Println(err.Error())
		return exists, &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return exists, nil
}

func insertUserDB(ctx context.Context, tx pgx.Tx, u *User) *utils.HTTPError {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}
	u.Password = string(bytes)

	if _, err := tx.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", u.Username, u.Password); err != nil {
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return nil
}

// come up with better name for this
func (us *UserStore) InsertUser(u *User) *utils.HTTPError {
	tx, err := us.Pool.Begin(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}
	defer tx.Rollback(context.TODO())

	exists, httpErr := checkUserExistsDB(context.TODO(), tx, u.Username)
	if httpErr != nil {
		return httpErr
	}

	if exists {
		return &utils.HTTPError{
			Code:    400,
			Message: "User already exists",
		}
	}

	httpErr = insertUserDB(context.TODO(), tx, u)
	if httpErr != nil {
		return httpErr
	}

	if err := tx.Commit(context.TODO()); err != nil {
		return &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return nil
}

func (us *UserStore) SelectUser(username string) (*User, *utils.HTTPError) {
	u := new(User)

	if err := us.Pool.QueryRow(context.TODO(), "SELECT username, password FROM users WHERE username = $1", username).Scan(&u.Username, &u.Password); err == pgx.ErrNoRows {
		return nil, &utils.HTTPError{
			Code:    400,
			Message: "User does not exist",
		}
	} else if err != nil {
		return nil, &utils.HTTPError{
			Code:    500,
			Message: "Internal server error",
		}
	}

	return u, nil
}
