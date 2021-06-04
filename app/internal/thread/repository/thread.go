package repository

import (
	"context"
	"database/sql"
	"strconv"

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
	var query string
	var queryParams []interface{}

	queryParams = append(queryParams,
		thread.Title,
		thread.Author,
		thread.Message,
		thread.Forum,
		thread.Slug,
	)
	if thread.Created != nil && *thread.Created != "" {
		query =
			`
		INSERT INTO threads (title, user_create, message, forum, slug, created) 
		VALUES ($1, $2, $3, $4, $5, $6) returning id
	`
		queryParams = append(queryParams, thread.Created)
	} else {
		query =
			`
		INSERT INTO threads (title, user_create, message, forum, slug) 
		VALUES ($1, $2, $3, $4, $5) returning id
	`
	}

	err = r.DB.QueryRow(query, queryParams...).Scan(&id)

	if err != nil {
		logger.Repo().AddFuncName("CreateThread").Error(ctx, err)
		return 0, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"thread id": id})
	return id, nil
}

func (r *repo) GetThreadBySlugOrId(ctx context.Context, slugOrId string) (*models.Thread, error) {

	thread := new(models.Thread)
	query :=
		`
		SELECT th.id, th.title, th.user_create, th.forum, 
		th.message, th.slug, th.created, sum(v.voice)
		FROM threads as th
		LEFT JOIN votes as v
		ON v.thread = th.id
	`

	if _, err := strconv.Atoi(slugOrId); err == nil {
		query += " WHERE th.id = $1"
	} else {
		query += " WHERE th.slug = $1"
	}

	query += " GROUP BY th.id"

	err := r.DB.QueryRow(query, slugOrId).Scan(
		&thread.Id,
		&thread.Title,
		&thread.Author,
		&thread.Forum,
		&thread.Message,
		&thread.Slug,
		&thread.Created,
		&thread.Votes,
	)
	if err == sql.ErrNoRows {
		logger.Repo().Info(ctx, logger.Fields{"thread": "not thread"})
		return nil, nil
	}
	if err != nil {
		logger.Repo().AddFuncName("GetThreadBySlugOrId").Error(ctx, err)
		return nil, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"thread": *thread})
	return thread, nil
}

func (r *repo) UpdateThreadBySlug(ctx context.Context, thread models.Thread) error {
	query :=
		`
		UPDATE threads SET title = $1, message = $2
		WHERE slug = $3
	`

	_, err := r.DB.Exec(query, thread.Title, thread.Message)
	if err != nil {
		logger.Repo().AddFuncName("UpdateThreadBySlug").Error(ctx, err)
		return err
	}

	return nil
}

func (r *repo) UpdateVote(ctx context.Context, vote models.Vote) error {
	query :=
		`
		UPDATE votes SET user_create = $1, thread = $2, voice = $3
		WHERE id = $4
	`

	_, err := r.DB.Exec(query, vote.User, vote.Thread, vote.Voice, vote.Id)
	if err != nil {
		logger.Repo().AddFuncName("UpdateVote").Error(ctx, err)
		return err
	}

	return nil
}

func (r *repo) CheckVote(ctx context.Context, vote models.Vote) (int, bool, error) {
	query :=
		`
		SELECT id
		FROM votes
		WHERE user_create = $1 AND thread = $2
	`

	id := new(int)
	err := r.DB.QueryRow(query, vote.User, vote.Thread).Scan(&id)

	if err == sql.ErrNoRows {
		logger.Repo().Info(ctx, logger.Fields{"vote": "not found"})
		return 0, false, nil
	}

	if err != nil {
		logger.Repo().AddFuncName("CheckVote").Error(ctx, err)
		return 0, false, err
	}

	return *id, true, nil
}

func (r *repo) AddVote(ctx context.Context, vote models.Vote) error {
	id := new(int)

	query :=
		`
		INSERT INTO votes (user_create, thread, voice) 
		VALUES ($1, $2, $3) returning id
	`
	err := r.DB.QueryRow(query, vote.User, vote.Thread, vote.Voice).Scan(&id)
	if err != nil {
		logger.Repo().AddFuncName("AddVote").Error(ctx, err)
		return err
	}

	return nil
}

func (r *repo) GetPosts(ctx context.Context, threadPosts models.ThreadPosts) ([]models.Post, error) {
	// TODO: подумать как здесь можно сделать покрасивее
	var queryParams []interface{}
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message,
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.thread = $1
	`

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Sort == "" {
		threadPosts.Sort = "flat"
	}
	switch threadPosts.Sort {
	case "tree":
		if threadPosts.Desc {
			if threadPosts.Since != "" {
				query += " AND p.id < $2"
				queryParams = append(queryParams, threadPosts.Since)
			}

			query += " ORDER BY p.tree[0] DESC, p.tree DESC"
		} else {
			if threadPosts.Since != "" {
				query += " AND p.id > $2"
				queryParams = append(queryParams, threadPosts.Since)
			}
			query += " ORDER BY p.tree"
		}

	case "parent_tree":
		if threadPosts.Desc {
			if threadPosts.Since != "" {
				query += " AND p.id < $2"
				queryParams = append(queryParams, threadPosts.Since)
			}

			query += " ORDER BY p.tree[1] DESC, p.tree ASC"
		} else {
			if threadPosts.Since != "" {
				query += " AND p.id > $2"
				queryParams = append(queryParams, threadPosts.Since)
			}
			query += " ORDER BY p.tree"
		}

	case "flat":
		if threadPosts.Desc {
			if threadPosts.Since != "" {
				query += " AND p.id < $2"
				queryParams = append(queryParams, threadPosts.Since)
			}

			query += " ORDER BY p.id DESC"
		} else {
			if threadPosts.Since != "" {
				query += " AND p.id > $2"
				queryParams = append(queryParams, threadPosts.Since)
			}
			query += " ORDER BY p.id"
		}
	}

	if threadPosts.Limit != "" {
		query += " LIMIT " + threadPosts.Limit
	}

	logger.Repo().Debug(ctx, logger.Fields{"query": query})

	threadsDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		logger.Repo().AddFuncName("GetPosts").Error(ctx, err)
		return nil, err
	}

	posts := make([]models.Post, 0)
	for threadsDB.Next() {
		post := new(models.Post)
		err := threadsDB.Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)

		if err != nil {
			logger.Repo().AddFuncName("GetPosts").Error(ctx, err)
			return nil, err
		}

		posts = append(posts, *post)
	}

	return posts, nil
}
