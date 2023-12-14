package users

import (
	"context"
	"database/sql"
	"open-ecm/apps"
	"time"

	"go.uber.org/zap"
)

type repository interface {
	getUser(context.Context, int64) (User, error)
	saveUser(context.Context, User) (User, error)
}

type repositoryImpl struct {
	db *sql.DB
}

func newRepository(db *sql.DB) repository {
	return repositoryImpl{db}
}

func (r repositoryImpl) saveUser(ctx context.Context, user User) (savedUser User, err error) {
	err = r.db.QueryRow(`
insert into users(username, email, password, created_at) values ($1, $2, $3, $4) returning id, username, email, password, created_at;
	`, user.Username, user.Email, user.Password, time.Now()).Scan(&savedUser.Id, &savedUser.Username, &savedUser.Email, &savedUser.Password, &savedUser.CreatedAt)
	return savedUser, err
}

func (r repositoryImpl) getUser(ctx context.Context, id int64) (user User, err error) {
	logger := apps.LoggerFromCtx(ctx)
	err = r.db.QueryRow(`
select id, username, email, password, created_at from users where id = $1;
	`, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		logger.Debug("user repository error", zap.String("error", err.Error()))
	}
	return
}
