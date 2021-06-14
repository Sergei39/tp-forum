package repository

import (
	"context"
	"strconv"

	forumModel "github.com/forums/app/internal/forum"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
	"github.com/jackc/pgx"
)

type repo struct {
	DB *pgx.ConnPool
}

func NewForumRepo(db *pgx.ConnPool) forumModel.ForumRepo {
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
		SELECT f.title, f.user_create, f.slug, f.threads, f.posts
		FROM forums as f
		WHERE f.slug = $1
	`
	err := r.DB.QueryRow(query, slug).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Threads,
		&forum.Posts,
	)

	if err == pgx.ErrNoRows {
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

func (r *repo) GetUsers(ctx context.Context, forumUsers models.ForumUsers) ([]models.User, error) {
	var queryParams []interface{}
	query :=
		`
		SELECT user_nickname, user_fullname, user_about, user_email
		FROM forums_users
		WHERE forum = $1
	`
	queryParams = append(queryParams, forumUsers.Slug)

	if forumUsers.Desc {
		if forumUsers.Since != "" {
			query += " AND user_nickname < $2"
			queryParams = append(queryParams, forumUsers.Since)
		}

		query += " ORDER BY user_nickname DESC"
	} else {
		if forumUsers.Since != "" {
			query += " AND user_nickname > $2"
			queryParams = append(queryParams, forumUsers.Since)
		}

		query += " ORDER BY user_nickname"
	}

	if forumUsers.Limit != "" {
		query += " LIMIT " + forumUsers.Limit
	}

	logger.Repo().AddFuncName("GetUsers").Debug(ctx, logger.Fields{"query": query})

	usersDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		logger.Repo().AddFuncName("GetUsers").Error(ctx, err)
		return nil, err
	}

	users := make([]models.User, 0)
	for usersDB.Next() {
		user := new(models.User)

		err := usersDB.Scan(
			&user.Nickname,
			&user.Fullname,
			&user.About,
			&user.Email,
		)

		if err != nil {
			logger.Repo().AddFuncName("GetUsers").Error(ctx, err)
			return nil, err
		}

		users = append(users, *user)
	}

	logger.Repo().Info(ctx, logger.Fields{"users": users})
	return users, nil
}

func (r *repo) GetThreads(ctx context.Context, forumThreads models.ForumThreads) ([]models.Thread, error) {
	// TODO: подумать как здесь можно сделать покрасивее
	var queryParams []interface{}
	query :=
		`
		SELECT th.id, th.title, th.user_create, th.forum, 
		th.message, th.slug, th.created, th.votes
		FROM threads as th
		WHERE th.forum = $1
	`
	queryParams = append(queryParams, forumThreads.Slug)

	if forumThreads.Desc {
		if forumThreads.Since != "" {
			query += " AND th.created <= $2"
			queryParams = append(queryParams, forumThreads.Since)
		}

		query += " ORDER BY th.created DESC"
	} else {
		if forumThreads.Since != "" {
			query += " AND th.created >= $2"
			queryParams = append(queryParams, forumThreads.Since)
		}
		query += " ORDER BY th.created"
	}

	if forumThreads.Limit != 0 {
		query += " LIMIT " + strconv.Itoa(forumThreads.Limit)
	}

	logger.Repo().Debug(ctx, logger.Fields{"query": query})

	threadsDB, err := r.DB.Query(query, queryParams...)
	if err != nil {
		logger.Repo().AddFuncName("GetThreads").Error(ctx, err)
		return nil, err
	}

	threads := make([]models.Thread, 0)
	for threadsDB.Next() {
		thread := new(models.Thread)

		err := threadsDB.Scan(
			&thread.Id,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Created,
			&thread.Votes,
		)

		if err != nil {
			logger.Repo().AddFuncName("GetThreads").Error(ctx, err)
			return nil, err
		}

		threads = append(threads, *thread)
	}

	logger.Repo().AddFuncName("GetThreads").Info(ctx, logger.Fields{"threads": threads})

	return threads, nil
}
