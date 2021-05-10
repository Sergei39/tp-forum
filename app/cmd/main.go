package main

import (
	"context"
	"database/sql"
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
	userRepository "github.com/forums/app/internal/user/repository"

	forumUsecase "github.com/forums/app/internal/forum/usecase"
	userUsecase "github.com/forums/app/internal/user/usecase"

	forumDelivery "github.com/forums/app/internal/forum/delivery"
	postDelivery "github.com/forums/app/internal/post/delivery"
	serviceDelivery "github.com/forums/app/internal/service/delivery"
	threadDelivery "github.com/forums/app/internal/thread/delivery"
	userDelivery "github.com/forums/app/internal/user/delivery"

	_ "github.com/lib/pq"
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
	userGroup := h.echo.Group("/user")
	userGroup.POST("/:nickname/create", h.user.CreateUser)
	userGroup.GET("/:nickname/profile", h.user.GetUser)
	userGroup.POST("/:nickname/profile", h.user.UpdateUser)

	forumGroup := h.echo.Group("/forum")
	forumGroup.POST("/create", h.forum.CreateForum)
	forumGroup.GET("/:slug/details", h.forum.GetDetails)
	forumGroup.POST("/:slug/create", h.forum.GetThreads)
	forumGroup.GET("/:slug/users", h.forum.GetUsers)
	forumGroup.GET("/:slug/threads", h.forum.GetThreads)

	postGroup := h.echo.Group("/post")
	postGroup.GET("/:id/details", h.post.GetDetails)
	postGroup.POST("/:id/details", h.post.UpdateDetails)

	serviceGroup := h.echo.Group("/service")
	serviceGroup.POST("/clear", h.service.ClearDb)
	serviceGroup.GET("/status", h.service.StatusDb)

	threadGroup := h.echo.Group("/thraed")
	threadGroup.POST("/:slag_or_id/create", h.thread.CreateThread)
	threadGroup.GET("/:slag_or_id/details", h.thread.GetDetails)
	threadGroup.POST("/:slag_or_id/details", h.thread.UpdateDetails)
	threadGroup.GET("/:slag_or_id/posts", h.thread.GetPosts)
	threadGroup.POST("/:slag_or_id/vote", h.thread.Vote)
}

func main() {
	logger.InitLogger()
	ctx := context.Background()

	e := echo.New()
	e.Use(custMiddleware.LogMiddleware)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s", config.DBUser, config.DBPass, config.DBName)
	db, err := sql.Open(config.PostgresDB, dsn)
	if err != nil {
		logger.Start().Fatal(ctx, err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)

	err = db.Ping()
	if err != nil {
		logger.Start().Fatal(ctx, err)
	}

	userRepo := userRepository.NewUserRepo(db)
	forumRepo := forumRepository.NewForumRepo(db)

	userUcase := userUsecase.NewUserUsecase(userRepo)
	forumUcase := forumUsecase.NewForumUsecase(forumRepo)

	userHandler := userDelivery.NewUserHandler(userUcase)
	forumHandler := forumDelivery.NewForumHandler(forumUcase)
	postHandler := postDelivery.NewPostHandler()
	serviceHandler := serviceDelivery.NewServiceHandler()
	threadHandler := threadDelivery.NewThreadHandler()

	handlers := Handler{
		echo:    e,
		user:    userHandler,
		forum:   forumHandler,
		post:    postHandler,
		service: serviceHandler,
		thread:  threadHandler,
	}

	router(handlers)

	logger.Start().Fatal(ctx, e.Start(":8080"))
}
