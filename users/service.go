package users

import (
	"context"
	"database/sql"
)

type service interface {
	createUser(context.Context, User) (User, error)
}

type serviceImpl struct {
	db   *sql.DB
	repo repository
}

func newService(db *sql.DB) service {
	return serviceImpl{db: db, repo: newRepository(db)}
}

func (s serviceImpl) createUser(ctx context.Context, user User) (savedUser User, err error) {
	return s.repo.persistUser(ctx, user)
}
