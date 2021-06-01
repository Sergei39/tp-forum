package repository

import (
	"context"
	"database/sql"

	postModel "github.com/forums/app/internal/post"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
)

type repo struct {
	DB *sql.DB
}

func NewPostRepo(db *sql.DB) postModel.PostRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) GetPost(ctx context.Context, id int) (*models.Post, error) {

	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message, 
		p.is_edited, f.slug, p.thread, p.created
		FROM posts as p
		JOIN forums as f
		ON p.forum = f.id
		WHERE p.id = $1
	`

	post := new(models.Post)
	err := r.DB.QueryRow(query, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Thread,
		&post.Created,
	)

	if err != nil {
		logger.Repo().AddFuncName("GetPost").Error(ctx, err)
		return nil, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"post": *post})
	return post, nil
}

func (r *repo) UpdateMessage(ctx context.Context, request models.MessagePostRequest) error {
	query :=
		`
		UPDATE post SET message = $1, is_edited = true
		WHERE id = $2
	`

	_, err := r.DB.Exec(query, request.Message, request.Id)
	if err != nil {
		logger.Repo().AddFuncName("UpdateMessage").Error(ctx, err)
		return err
	}

	return nil
}
