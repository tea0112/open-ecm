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
	deleteUser(context.Context, int64) (*User, error)
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
	from users where id = $1 and deleted_at is null;
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
		select id, username, email, password, created_at, updated_at from users
		where deleted_at is null;
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

func (r repositoryImpl) deleteUser(ctx context.Context, id int64) (*User, error) {
	row := r.db.QueryRow(`
		update users
		set deleted_at = $1
		where id = $2 and deleted_at is null
		returning id, username, email, password, created_at, updated_at;
	`, time.Now(), id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var deletedUser User
	err := row.Scan(
		&deletedUser.Id,
		&deletedUser.Username,
		&deletedUser.Email,
		&deletedUser.Password,
		&deletedUser.CreatedAt,
		&deletedUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &deletedUser, nil
}
