package repository

import (
	"context"
	"database/sql"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
)

type repo struct {
	DB *sql.DB
}

func NewThreadRepo(db *sql.DB) threadModel.ThreadRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) CreateThread(ctx context.Context, thread models.Thread) (id int, err error) {

	query :=
		`
		INSERT INTO threads (title, user_create, message, created) 
		VALUES ($1, $2, $3, $4) returning id
	`
	err = r.DB.QueryRow(query,
		thread.Title,
		thread.Author,
		thread.Message,
		thread.Created).Scan(&id)

	if err != nil {
		logger.Repo().AddFuncName("CreateThread").Error(ctx, err)
		return 0, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"thread id": id})
	return id, nil
}

// func (r *repo) GetThreadBySlug(ctx context.Context, slug string) (*models.Forum, error) {

// 	forum := new(models.Forum)
// 	query :=
// 		`
// 		SELECT f.title, f.user_create, f.slug, count(p), count(th)
// 		FROM forums as f
// 		JOIN posts as p
// 		ON p.forum = f.id
// 		JOIN threads as th
// 		ON th.forum = f.id
// 		WHERE f.slug = $1
// 	`
// 	err := r.DB.QueryRow(query, slug).Scan(
// 		&forum.Title,
// 		&forum.User,
// 		&forum.Slug,
// 		&forum.Posts,
// 		&forum.Threads,
// 	)
// 	if err == sql.ErrNoRows {
// 		logger.Repo().Info(ctx, logger.Fields{"forum": "not forum"})
// 		return nil, nil
// 	}
// 	if err != nil {
// 		logger.Repo().AddFuncName("GetForumBySlug").Error(ctx, err)
// 		return nil, err
// 	}

// 	logger.Repo().Debug(ctx, logger.Fields{"forum": *forum})
// 	return forum, nil
// }
