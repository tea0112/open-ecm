package roles

import (
	"context"
	"database/sql"
	"open-ecm/apps"
	"time"
)

type repository interface {
	saveRole(context.Context, *sql.Tx, Role) (*Role, error)
	getRole(context.Context, *sql.Tx, int64) (*Role, error)
	getRoles(context.Context, *sql.Tx) ([]Role, error)
	deleteRoles(context.Context, *sql.Tx, Role) (*Role, error)
}

type repositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository {
	return repositoryImpl{db: db}
}

func (r repositoryImpl) saveRole(ctx context.Context, tx *sql.Tx, role Role) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	row := tx.QueryRow(`
		insert into (name, created_at)
		values ($1, $2)
		returning name, created_at;
	`, role.Name, time.Now())
	err := row.Err()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var createdRole Role
	err = row.Scan(&createdRole.Id, &createdRole.Name, &createdRole.CreatedAt)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &createdRole, nil
}

func (r repositoryImpl) getRole(ctx context.Context, tx *sql.Tx, id int64) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	row := tx.QueryRow(`
		select id, name, created_at, updated_at
		from roles
		where deleted_at is null and id = $1;
	`, id)
	err := row.Err()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var retrievedRole Role
	err = row.Scan(&retrievedRole.Id, &retrievedRole.Name, &retrievedRole.CreatedAt, &retrievedRole.UpdatedAt)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &retrievedRole, nil
}

func (r repositoryImpl) getRoles(ctx context.Context, tx *sql.Tx) ([]Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	rows, err := tx.Query(`
		select id, name, created_at, updated_at
		from roles
		where deleted_at is null;
	`)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var retrievedRoles []Role
	for rows.Next() {
		var retrievedRole Role
		err = rows.Scan(&retrievedRole.Id, &retrievedRole.Name, &retrievedRole.CreatedAt, &retrievedRole.UpdatedAt)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		retrievedRoles = append(retrievedRoles, retrievedRole)
	}

	return retrievedRoles, nil
}

func (r repositoryImpl) deleteRoles(ctx context.Context, tx *sql.Tx, role Role) (*Role, error) {
	logger := apps.LoggerFromCtx(ctx)

	row := tx.QueryRow(`
		delete from roles
		where id = $1 and deleted_at is null
		returning id, name, created_at, updated_at, deleted_at;
	`, role.Id)
	err := row.Err()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	var deletedRole Role
	err = row.Scan(&role.Id, &role.Name, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return &deletedRole, nil
}
