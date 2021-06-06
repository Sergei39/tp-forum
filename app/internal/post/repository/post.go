package repository

import (
	"context"
	"fmt"

	postModel "github.com/forums/app/internal/post"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewPostRepo(db *pgx.ConnPool) postModel.PostRepo {
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
		ON p.forum = f.slug
		WHERE p.id = $1
	`

	post := new(models.Post)
	err := r.DB.QueryRow(query, id).Scan(
		&post.Id,
		&post.Parent,
		&post.Author,
		&post.Message,
		&post.IsEdited,
		&post.Forum,
		&post.Thread,
		&post.Created,
	)

	if err == pgx.ErrNoRows {
		logger.Repo().Info(ctx, logger.Fields{"post": "not post"})
		return nil, nil
	}

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
		UPDATE posts SET message = $1, is_edited = true
		WHERE id = $2
	`

	_, err := r.DB.Exec(query, request.Message, request.Id)
	if err != nil {
		logger.Repo().AddFuncName("UpdateMessage").Error(ctx, err)
		return err
	}

	return nil
}

func (r *repo) CreatePosts(ctx context.Context, posts []models.Post) ([]models.Post, error) {
	var queryParams []interface{}
	query := "INSERT INTO posts (parent, user_create, message, forum, thread, created, tree) VALUES "

	for i, post := range posts {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, (SELECT tree FROM posts WHERE id = $%d) || ARRAY[nextval('post_tree_id')])",
			i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6, i*6+1)

		if i != len(posts)-1 {
			query += ","
		}

		queryParams = append(queryParams, post.Parent, post.Author, post.Message, post.Forum, post.Thread, post.Created)
	}

	query += " returning id, created"

	logger.Repo().AddFuncName("CreatePosts").Debug(ctx, logger.Fields{"query": query})

	postsDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		logger.Repo().AddFuncName("CreatePosts_Query").Info(ctx, logger.Fields{"Error": err})
		return nil, err
	}

	i := 0
	for postsDB.Next() {
		err = postsDB.Scan(
			&(posts[i].Id),
			&(posts[i].Created),
		)

		if err != nil {
			logger.Repo().AddFuncName("CreatePosts_Scan").Error(ctx, err)
			return nil, err
		}
		i++
	}

	if dbErr, ok := postsDB.Err().(pgx.PgError); ok {
		return nil, dbErr
	}

	return posts, nil
}

func (r *repo) CreateForumsUsers(ctx context.Context, posts []models.Post) error {
	var queryParams []interface{}
	query := "INSERT INTO forums_users (forum, user_create) VALUES "

	for i, post := range posts {
		query += fmt.Sprintf("($%d, $%d)",
			i*2+1, i*2+2)

		if i != len(posts)-1 {
			query += ","
		}

		queryParams = append(queryParams, post.Forum, post.Author)
	}

	query += " ON CONFLICT DO NOTHING"

	logger.Repo().AddFuncName("CreateForumsUsers").Debug(ctx, logger.Fields{"query": query})

	_, err := r.DB.Exec(query, queryParams...)
	if err != nil {
		logger.Repo().AddFuncName("CreateForumsUsers").Error(ctx, err)
		return err
	}

	return nil
}
