package repository

import (
	"context"

	serviceModel "github.com/forums/app/internal/service"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewServiceRepo(db *pgx.ConnPool) serviceModel.ServiceRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) ClearDb(ctx context.Context) error {

	query :=
		`
		TRUNCATE users, forums, threads, posts, forums_users, votes CASCADE
	`
	result, err := r.DB.Exec(query)
	if err != nil {
		logger.Repo().AddFuncName("ClearDb").Error(ctx, err)
		return err
	}

	logger.Repo().Info(ctx, logger.Fields{"result": result})
	return nil
}

func (r *repo) StatusDb(ctx context.Context) (*models.InfoStatus, error) {
	// TODO: возможно сделать не прямой запрос кол-ва а через другие средства
	info := new(models.InfoStatus)
	var err error

	info.User, err = r.getUsersNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Forum, err = r.getForumsNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Thread, err = r.getThreadsNumber(ctx)
	if err != nil {
		return nil, err
	}

	info.Post, err = r.getPostsNumber(ctx)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (r *repo) getUsersNumber(ctx context.Context) (number int, err error) {
	query :=
		`
		SELECT COUNT(*) FROM users
	`

	err = r.DB.QueryRow(query).Scan(&number)
	if err != nil {
		logger.Repo().AddFuncName("getUsersNumber").Error(ctx, err)
		return 0, err
	}

	return number, nil
}

func (r *repo) getForumsNumber(ctx context.Context) (number int, err error) {
	query :=
		`
		SELECT COUNT(*) FROM forums
	`

	err = r.DB.QueryRow(query).Scan(&number)
	if err != nil {
		logger.Repo().AddFuncName("getForumsNumber").Error(ctx, err)
		return 0, err
	}

	return number, nil
}

func (r *repo) getThreadsNumber(ctx context.Context) (number int, err error) {
	query :=
		`
		SELECT COUNT(*) FROM threads
	`

	err = r.DB.QueryRow(query).Scan(&number)
	if err != nil {
		logger.Repo().AddFuncName("getThreadsNumber").Error(ctx, err)
		return 0, err
	}

	return number, nil
}

func (r *repo) getPostsNumber(ctx context.Context) (number int, err error) {
	query :=
		`
		SELECT COUNT(*) FROM posts
	`

	err = r.DB.QueryRow(query).Scan(&number)
	if err != nil {
		logger.Repo().AddFuncName("getPostsNumber").Error(ctx, err)
		return 0, err
	}

	return number, nil
}
