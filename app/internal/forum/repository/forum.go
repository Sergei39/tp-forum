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

func NewForumRepo(db *sql.DB) forumModel.ForumRepo {
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
		SELECT f.title, f.user_create, f.slug, count(p), count(th)
		FROM forums as f
		JOIN posts as p 
		ON p.forum = f.id
		JOIN threads as th
		ON th.forum = f.id
		WHERE f.slug = $1
	`
	err := r.DB.QueryRow(query, slug).Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
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

func (r *repo) GetUsers(ctx context.Context, forumUsers models.ForumUsers) ([]models.User, error) {
	// TODO: доделать правильный запрос
	query :=
		`
		SELECT DISTINCT u.nickname, u.fullname, u.about, u.email
		FROM forum as f
		JOIN threads as th 
		ON th.forum = f.id
		JOIN posts as p
		ON p.forum = f.id
		JOIN users as u
		ON u.nickname = th.user_create OR u.nickname = p.user_create
		WHERE f.slug = $1
		ORDER BY u.nickname
	`
	usersDB, err := r.DB.Query(query, forumUsers.Slug)
	if err != nil {
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

	return users, nil
}

func (r *repo) GetThreads(ctx context.Context, forumThreads models.ForumThreads) ([]models.Thread, error) {
	// TODO: доделать правильный запрос
	query :=
		`
		SELECT DISTINCT th.id, th.title, th.user_create, f.slug, 
		th.message, count(v), th.slug, th.created
		FROM forum as f
		JOIN threads as th 
		ON th.forum = f.id
		JOIN votes as v
		ON v.thread = th.id
		WHERE f.slug = $1
		ORDER BY th.created
	`
	threadsDB, err := r.DB.Query(query, forumThreads.Slug)
	if err != nil {
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
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)

		if err != nil {
			logger.Repo().AddFuncName("GetThreads").Error(ctx, err)
			return nil, err
		}

		threads = append(threads, *thread)
	}

	return threads, nil
}
