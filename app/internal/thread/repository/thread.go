package repository

import (
	"context"
	"strconv"

	"github.com/jackc/pgx"

	threadModel "github.com/forums/app/internal/thread"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewThreadRepo(db *pgx.ConnPool) threadModel.ThreadRepo {
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
	if thread.Created != nil {
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
		th.message, th.slug, th.created, th.votes
		FROM threads as th
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
	if err == pgx.ErrNoRows {
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

	_, err := r.DB.Exec(query, thread.Title, thread.Message, thread.Slug)
	if err != nil {
		logger.Repo().AddFuncName("UpdateThreadBySlug").Error(ctx, err)
		return err
	}

	return nil
}

func (r *repo) UpdateVote(ctx context.Context, vote models.Vote) error {
	query :=
		`
		UPDATE votes SET voice = $1
		WHERE user_create = $2 AND thread = $3
	`

	_, err := r.DB.Exec(query, vote.Voice, vote.User, vote.Thread)
	if err != nil {
		logger.Repo().AddFuncName("UpdateVote").Error(ctx, err)
		return err
	}

	return nil
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
		return err
	}

	logger.Repo().AddFuncName("AddVote").Info(ctx, logger.Fields{"vote id": id})
	return nil
}

func (r *repo) treeSort(ctx context.Context, threadPosts models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message,
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.thread = $1
	`

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Desc {
		if threadPosts.Since != "" {
			query += " AND p.tree < (SELECT p2.tree from posts AS p2 WHERE p2.id = $2)"
			queryParams = append(queryParams, threadPosts.Since)
		}

		query += " ORDER BY p.tree[0] DESC, p.tree DESC"
	} else {
		if threadPosts.Since != "" {
			query += " AND p.tree > (SELECT p2.tree from posts AS p2 WHERE p2.id = $2)"
			queryParams = append(queryParams, threadPosts.Since)
		}
		query += " ORDER BY p.tree"
	}

	if threadPosts.Limit != "" {
		query += " LIMIT " + threadPosts.Limit
	}

	return query, queryParams
}

const selectParentTreeLimitAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.tree[1] IN (
		SELECT p2.tree[1]
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent = 0
		ORDER BY p2.tree
		LIMIT $2
	)
	ORDER BY p.tree
`

const selectParentTreeLimitDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.tree[1] IN (
		SELECT p2.tree[1]
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent = 0
		ORDER BY p2.tree DESC
		LIMIT $2
	)
	ORDER BY p.tree[1] DESC, p.tree ASC
`

const selectParentTreeSinceLimitAsc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.tree[1] IN (
		SELECT p2.tree[1]
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent = 0 AND p2.tree[1] > (SELECT p3.tree[1] from posts p3 where p3.id = $2)
		ORDER BY p2.tree
		LIMIT $3
	)
	ORDER BY p.tree
`

const selectParentTreeSinceLimitDesc = `
	SELECT p.id, p.parent, p.user_create, p.message,
	p.is_edited, p.forum, p.thread, p.created
	FROM posts as p
	WHERE p.thread = $1 and p.tree[1] IN (
		SELECT p2.tree[1]
		FROM posts p2
		WHERE p2.thread = $1 AND p2.parent = 0 AND p2.tree[1] < (SELECT p3.tree[1] from posts p3 where p3.id = $2)
		ORDER BY p2.tree DESC
		LIMIT $3
	)
	ORDER BY p.tree[1] DESC, p.tree ASC
`

func (r *repo) parentTreeSort(ctx context.Context, threadPosts models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	var query string

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Desc {
		if threadPosts.Since != "" {
			query = selectParentTreeSinceLimitDesc
			queryParams = append(queryParams, threadPosts.Since)
		} else {
			query = selectParentTreeLimitDesc
		}
	} else {
		if threadPosts.Since != "" {
			query = selectParentTreeSinceLimitAsc
			queryParams = append(queryParams, threadPosts.Since)
		} else {
			query = selectParentTreeLimitAsc
		}
	}

	if threadPosts.Limit != "" {
		queryParams = append(queryParams, threadPosts.Limit)
	} else {
		queryParams = append(queryParams, "ALL")
	}

	return query, queryParams
}

func (r *repo) flatSort(ctx context.Context, threadPosts models.ThreadPosts) (string, []interface{}) {
	var queryParams []interface{}
	query :=
		`
		SELECT p.id, p.parent, p.user_create, p.message,
		p.is_edited, p.forum, p.thread, p.created
		FROM posts as p
		WHERE p.thread = $1
	`

	queryParams = append(queryParams, threadPosts.ThreadId)

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

	if threadPosts.Limit != "" {
		query += " LIMIT " + threadPosts.Limit
	}

	return query, queryParams
}

func (r *repo) GetPosts(ctx context.Context, threadPosts models.ThreadPosts) ([]models.Post, error) {
	// TODO: подумать как здесь можно сделать покрасивее
	var queryParams []interface{}
	var query string

	queryParams = append(queryParams, threadPosts.ThreadId)

	if threadPosts.Sort == "" {
		threadPosts.Sort = "flat"
	}
	switch threadPosts.Sort {
	case "tree":
		query, queryParams = r.treeSort(ctx, threadPosts)

	case "parent_tree":
		query, queryParams = r.parentTreeSort(ctx, threadPosts)

	case "flat":
		query, queryParams = r.flatSort(ctx, threadPosts)
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
