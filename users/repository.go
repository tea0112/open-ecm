package users

import (
	"context"
	"database/sql"
	"time"
)

type repository interface {
	persistUser(context.Context, User) (User, error)
}

type repositoryImpl struct {
	db *sql.DB
}

func newRepository(db *sql.DB) repository {
	return repositoryImpl{db}
}

func (r repositoryImpl) persistUser(ctx context.Context, user User) (savedUser User, err error) {
	err = r.db.QueryRow(`
insert into users(username, email, password, created_at) values ($1, $2, $3, $4) returning id, username, email, password, created_at;
	`, user.Username, user.Email, user.Password, time.Now()).Scan(&savedUser.Id, &savedUser.Username, &savedUser.Email, &savedUser.Password, &savedUser.CreatedAt)
	return savedUser, err
}
