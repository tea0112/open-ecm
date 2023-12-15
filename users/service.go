package users

import (
	"context"
	"database/sql"
)

type service interface {
	saveUser(context.Context, User) (*User, error)
	getUser(context.Context, int64) (*User, error)
	getUsers(context.Context) ([]User, error)
}

type serviceImpl struct {
	db   *sql.DB
	repo repository
}

func newService(db *sql.DB) service {
	return serviceImpl{db: db, repo: newRepository(db)}
}

func (s serviceImpl) saveUser(ctx context.Context, user User) (*User, error) {
	savedUser, err := s.repo.saveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

func (s serviceImpl) getUser(ctx context.Context, id int64) (*User, error) {
	retrievedUser, err := s.repo.getUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return retrievedUser, nil
}

func (s serviceImpl) getUsers(ctx context.Context) ([]User, error) {
	retrievedUsers, err := s.repo.getUsers(ctx)
	if err != nil {
		return nil, err
	}
	return retrievedUsers, nil
}
