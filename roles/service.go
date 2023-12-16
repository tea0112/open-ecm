package roles

import (
	"context"
	"database/sql"
	"fmt"
	"open-ecm/apps"
)

type service interface {
	saveRole(context.Context, Role) (*Role, error)
	getRole(context.Context, int64) (*Role, error)
	getRoles(context.Context) ([]Role, error)
	deleteRole(context.Context, int64) (*Role, error)
}

type serviceImpl struct {
	db *sql.DB
	r repository
}

func NewService(db *sql.DB) service {
	return serviceImpl{db: db, r: NewRepository(db)}
}

func (s serviceImpl) saveRole(ctx context.Context, newRole Role) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	savedRole, err := s.r.saveRole(ctx, tx, newRole)
	if err != nil {
		logger.Error(err.Error())

		rbErr := tx.Rollback()
		if rbErr != nil {
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return savedRole, nil
}

func (s serviceImpl) getRole(ctx context.Context, id int64) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
	}

	retrievedRole, err := s.r.getRole(ctx, tx, id)
	if err != nil {
		logger.Error(err.Error())
		rbErr := tx.Rollback()
		if rbErr != nil {
			logger.Error(rbErr.Error())
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return retrievedRole, nil
} 

func (s serviceImpl) getRoles(ctx context.Context) ([]Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	retrievedRoles, err := s.r.getRoles(ctx, tx)
	if err != nil {
		logger.Error(err.Error())
		rbErr := tx.Rollback()
		if rbErr != nil {
			logger.Error(rbErr.Error())
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return retrievedRoles, nil
}

func (s serviceImpl) deleteRole(ctx context.Context, id int64) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	retrievedRole, err := s.r.getRole(ctx, tx, id)
	if err != nil {
		logger.Error(err.Error())
		rbErr := tx.Rollback()
		if rbErr != nil {
			logger.Error(rbErr.Error())
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}

	deletedRole, err := s.r.deleteRoles(ctx, tx, *retrievedRole)
	if err != nil {
		logger.Error(err.Error())
		rbErr := tx.Rollback()
		if rbErr != nil {
			logger.Error(rbErr.Error())
			return nil, fmt.Errorf("%s: %w", err.Error(), rbErr)
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return deletedRole, nil
}