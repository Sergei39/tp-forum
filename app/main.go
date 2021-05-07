package main

import (
	"github.com/labstack/echo/v4"

	forumModels "github.com/forums/internal/forum"
	custMiddleware "github.com/forums/internal/middleware"
	postModels "github.com/forums/internal/post"
	serviceModels "github.com/forums/internal/service"
	threadModels "github.com/forums/internal/thread"
	userModels "github.com/forums/internal/user"
	"github.com/forums/utils/logger"

	forumDelivery "github.com/forums/internal/forum/delivery"
	postDelivery "github.com/forums/internal/post/delivery"
	serviceDelivery "github.com/forums/internal/service/delivery"
	threadDelivery "github.com/forums/internal/thread/delivery"
	userDelivery "github.com/forums/internal/user/delivery"
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
	userGroup.POST("/:username/create", h.user.CreateUser)
	userGroup.GET("/:username/profile", h.user.GetUser)
	userGroup.POST("/:username/profile", h.user.UpdateUser)

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

	e := echo.New()
	e.Use(custMiddleware.LogMiddleware)

	userHandler := userDelivery.NewUserHandler()
	forumHandler := forumDelivery.NewForumHandler()
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

	e.Logger.Fatal(e.Start(":8080"))
}
