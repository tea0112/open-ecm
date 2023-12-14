package users

import (
	"context"
	"database/sql"
)

type service interface {
	saveUser(context.Context, User) (User, error)
	getUser(context.Context, int64) (User, error)
}

type serviceImpl struct {
	db   *sql.DB
	repo repository
}

func newService(db *sql.DB) service {
	return serviceImpl{db: db, repo: newRepository(db)}
}

func (s serviceImpl) saveUser(ctx context.Context, user User) (savedUser User, err error) {
	return s.repo.saveUser(ctx, user)
}

func (s serviceImpl) getUser(ctx context.Context, id int64) (user User, err error) {
	return s.repo.getUser(ctx, id)
}
