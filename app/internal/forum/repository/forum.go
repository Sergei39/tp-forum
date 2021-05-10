package repository

import (
	"context"
	"database/sql"

	forumModel "github.com/forums/app/internal/forum"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
)

type repo struct {
	DB *sql.DB
}

func NewForumRepo(db *sql.DB) forumModel.UserRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateForum(ctx context.Context, forum models.Forum) (id int, err error) {

	query :=
		`
		INSERT INTO forums (title, user_create, slug) 
		VALUES ($1, $2, $3) returning id
	`
	err = r.DB.QueryRow(query,
		forum.Title,
		forum.User,
		forum.Slug).Scan(&id)

	if err != nil {
		logger.Repo().AddFuncName("CreateForum").Error(ctx, err)
		return 0, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"forum id": id})
	return id, nil
}

func (r *repo) GetForumBySlug(ctx context.Context, slug string) (*models.Forum, error) {

	forum := new(models.Forum)
	query :=
		`
		SELECT title, user_create, slug
		FROM forums WHERE slug = $1
	`
	err := r.DB.QueryRow(query, slug).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
	)
	if err == sql.ErrNoRows {
		logger.Repo().Info(ctx, logger.Fields{"forum": "not forum"})
		return nil, nil
	}
	if err != nil {
		logger.Repo().AddFuncName("GetForumBySlug").Error(ctx, err)
		return nil, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"forum": *forum})
	return forum, nil
}
