package repository

import (
	"context"
	"database/sql"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/logger"
)

type repo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) userModel.UserRepo {
	return &repo{
		DB: db,
	}
}

func (r *repo) GetUserByName(ctx context.Context, name string) (
	*models.User, error) {

	user := new(models.User)
	query :=
		`
		SELECT nickname, fullname, about, email
		FROM users WHERE nickname = $1
	`
	err := r.DB.QueryRow(query, name).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		logger.Repo().Info(ctx, logger.Fields{"user": "not user"})
		return nil, nil
	}
	if err != nil {
		logger.Repo().Error(ctx, err)
		return nil, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"user": *user})
	return user, nil
}

func (r *repo) CreateUser(ctx context.Context, user models.User) (
	id int, err error) {

	query :=
		`
		INSERT INTO users (nickname, fullname, about, email) 
		VALUES ($1, $2, $3, $4) returning id

	`
	err = r.DB.QueryRow(query,
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email).Scan(&id)

	if err != nil {
		logger.Repo().AddFuncName("CreateUser").Error(ctx, err)
		return 0, err
	}

	logger.Repo().Debug(ctx, logger.Fields{"user id": id})
	return id, nil
}
