package users

import (
	"context"
	"database/sql"
	"fmt"
	"open-ecm/apps"
)

type service interface {
	saveUser(context.Context, User) (*User, error)
	getUser(context.Context, int64) (*User, error)
	getUsers(context.Context) ([]User, error)
	deleteUser(context.Context, int64) (*User, error)
}

type serviceImpl struct {
	db   *sql.DB
	repo repository
}

func newService(db *sql.DB) service {
	return serviceImpl{db: db, repo: newRepository(db)}
}

func (s serviceImpl) saveUser(ctx context.Context, user User) (*User, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	savedUser, err := s.repo.saveUser(ctx, tx, user)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return savedUser, nil
}

func (s serviceImpl) getUser(ctx context.Context, id int64) (*User, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	retrievedUser, err := s.repo.getUser(ctx, tx, id)
	if err != nil {
		rolErr := tx.Rollback()
		if rolErr != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), rolErr)
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return retrievedUser, nil
}

func (s serviceImpl) getUsers(ctx context.Context) ([]User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	retrievedUsers, err := s.repo.getUsers(ctx, tx)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return retrievedUsers, nil
}

func (s serviceImpl) deleteUser(ctx context.Context, id int64) (*User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	deletedUser, err := s.repo.deleteUser(ctx, tx, id)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return deletedUser, nil
}
