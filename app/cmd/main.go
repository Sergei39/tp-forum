package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/forums/app/config"
	forumModels "github.com/forums/app/internal/forum"
	postModels "github.com/forums/app/internal/post"
	serviceModels "github.com/forums/app/internal/service"
	threadModels "github.com/forums/app/internal/thread"
	userModels "github.com/forums/app/internal/user"
	custMiddleware "github.com/forums/app/middleware"
	"github.com/forums/utils/logger"

	forumRepository "github.com/forums/app/internal/forum/repository"
	postRepository "github.com/forums/app/internal/post/repository"
	serviceRepository "github.com/forums/app/internal/service/repository"
	threadRepository "github.com/forums/app/internal/thread/repository"
	userRepository "github.com/forums/app/internal/user/repository"

	forumUsecase "github.com/forums/app/internal/forum/usecase"
	postUsecase "github.com/forums/app/internal/post/usecase"
	serviceUsecase "github.com/forums/app/internal/service/usecase"
	threadUsecase "github.com/forums/app/internal/thread/usecase"
	userUsecase "github.com/forums/app/internal/user/usecase"

	forumDelivery "github.com/forums/app/internal/forum/delivery"
	postDelivery "github.com/forums/app/internal/post/delivery"
	serviceDelivery "github.com/forums/app/internal/service/delivery"
	threadDelivery "github.com/forums/app/internal/thread/delivery"
	userDelivery "github.com/forums/app/internal/user/delivery"

	"github.com/jackc/pgx"
)

type Handler struct {
	echo    *echo.Echo
	user    userModels.UserHandler
	forum   forumModels.ForumHandler
	post    postModels.PostHandler
	service serviceModels.ServiceHandler
	thread  threadModels.ThreadHandler
}

func router(h Handler) {
	apiGroup := h.echo.Group("/api")
	userGroup := apiGroup.Group("/user")
	userGroup.POST("/:nickname/create", h.user.CreateUser)
	userGroup.GET("/:nickname/profile", h.user.GetUser)
	userGroup.POST("/:nickname/profile", h.user.UpdateUser)

	forumGroup := apiGroup.Group("/forum")
	forumGroup.POST("/create", h.forum.CreateForum)
	forumGroup.GET("/:slug/details", h.forum.GetDetails)
	forumGroup.POST("/:slug/create", h.thread.CreateThread)
	forumGroup.GET("/:slug/users", h.forum.GetUsers)
	forumGroup.GET("/:slug/threads", h.forum.GetThreads)

	postGroup := apiGroup.Group("/post")
	postGroup.GET("/:id/details", h.post.GetDetails)
	postGroup.POST("/:id/details", h.post.UpdateDetails)

	serviceGroup := apiGroup.Group("/service")
	serviceGroup.POST("/clear", h.service.ClearDb)
	serviceGroup.GET("/status", h.service.StatusDb)

	threadGroup := apiGroup.Group("/thread")
	threadGroup.POST("/:slug_or_id/create", h.post.CreatePosts)
	threadGroup.GET("/:slug_or_id/details", h.thread.GetDetails)
	threadGroup.POST("/:slug_or_id/details", h.thread.UpdateDetails)
	threadGroup.GET("/:slug_or_id/posts", h.thread.GetPosts)
	threadGroup.POST("/:slug_or_id/vote", h.thread.Vote)
}

func main() {
	logger.InitLogger()
	ctx := context.Background()

	e := echo.New()
	e.Use(custMiddleware.LogMiddleware)

	connectionString := "postgres://" + config.DBUser + ":" + config.DBPass +
		"@localhost/" + config.DBName + "?sslmode=disable"

	configDB, err := pgx.ParseURI(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     configDB,
			MaxConnections: 2000,
		})

	if err != nil {
		fmt.Println(err)
		return
	}

	userRepo := userRepository.NewUserRepo(db)
	forumRepo := forumRepository.NewForumRepo(db)
	serviceRepo := serviceRepository.NewServiceRepo(db)
	postRepo := postRepository.NewPostRepo(db)
	threadRepo := threadRepository.NewThreadRepo(db)

	userUcase := userUsecase.NewUserUsecase(userRepo)
	forumUcase := forumUsecase.NewForumUsecase(forumRepo, userRepo)
	serviceUcase := serviceUsecase.NewServiceUsecase(serviceRepo)
	postUcase := postUsecase.NewPostUsecase(postRepo, userRepo, threadRepo, forumRepo)
	threadUcase := threadUsecase.NewThreadUsecase(threadRepo, userRepo, forumRepo)

	userHandler := userDelivery.NewUserHandler(userUcase)
	forumHandler := forumDelivery.NewForumHandler(forumUcase)
	postHandler := postDelivery.NewPostHandler(postUcase)
	serviceHandler := serviceDelivery.NewServiceHandler(serviceUcase)
	threadHandler := threadDelivery.NewThreadHandler(threadUcase)

	handlers := Handler{
		echo:    e,
		user:    userHandler,
		forum:   forumHandler,
		post:    postHandler,
		service: serviceHandler,
		thread:  threadHandler,
	}

	router(handlers)

	logger.Start().Fatal(ctx, e.Start(":5000"))
}
