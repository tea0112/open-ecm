package users

import (
	"context"
	"database/sql"
	"open-ecm/apps"
	"time"

	"go.uber.org/zap"
)

type repository interface {
	getUser(context.Context, int64) (*User, error)
	getUsers(context.Context) ([]User, error)
	saveUser(context.Context, User) (*User, error)
}

type repositoryImpl struct {
	db *sql.DB
}

func newRepository(db *sql.DB) repository {
	return repositoryImpl{db}
}

func (r repositoryImpl) saveUser(ctx context.Context, user User) (*User, error) {
	var savedUser User
	err := r.db.QueryRow(`
	insert into users(username, email, password, created_at)
	values ($1, $2, $3, $4)
	returning id, username, email, password, created_at;
	`, user.Username, user.Email, user.Password, time.Now()).
		Scan(&savedUser.Id, &savedUser.Username, &savedUser.Email, &savedUser.Password, &savedUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &savedUser, nil
}

func (r repositoryImpl) getUser(ctx context.Context, id int64) (*User, error) {
	var retrievedUser User
	logger := apps.LoggerFromCtx(ctx)
	err := r.db.QueryRow(`
	select id, username, email, password, created_at
	from users where id = $1;
	`, id).Scan(&retrievedUser.Id, &retrievedUser.Username, &retrievedUser.Email, &retrievedUser.Password, &retrievedUser.CreatedAt)
	if err != nil {
		logger.Debug("user repository error", zap.String("error", err.Error()))
	}
	if err != nil {
		return nil, err
	}
	return &retrievedUser, nil
}

func (r repositoryImpl) getUsers(ctx context.Context) ([]User, error) {
	var retrievedUsers []User

	rows, err := r.db.Query(`
	select id, username, email, password, created_at, updated_at from users;
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		retrievedUsers = append(retrievedUsers, user)
	}

	return retrievedUsers, nil
}
