package main

import (
	"errors"
	"fmt"
	"net/http"

	custMiddleware "github.com/forums/app/middleware"

	"github.com/forums/app/config"
	forumModels "github.com/forums/app/internal/forum"
	postModels "github.com/forums/app/internal/post"
	serviceModels "github.com/forums/app/internal/service"
	threadModels "github.com/forums/app/internal/thread"
	userModels "github.com/forums/app/internal/user"
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

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type Handler struct {
	user    userModels.UserHandler
	forum   forumModels.ForumHandler
	post    postModels.PostHandler
	service serviceModels.ServiceHandler
	thread  threadModels.ThreadHandler
}

func newRouter(h Handler) *mux.Router {
	router := mux.NewRouter()
	router.Use(custMiddleware.LogMiddleware)

	user := router.PathPrefix("/api/user").Subrouter()
	user.HandleFunc("/{nickname}/create", h.user.CreateUser).Methods(http.MethodPost)
	user.HandleFunc("/{nickname}/profile", h.user.GetUser).Methods(http.MethodGet)
	user.HandleFunc("/{nickname}/profile", h.user.UpdateUser).Methods(http.MethodPost)

	forum := router.PathPrefix("/api/forum").Subrouter()
	forum.HandleFunc("/create", h.forum.CreateForum).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/details", h.forum.GetDetails).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/create", h.thread.CreateThread).Methods(http.MethodPost)
	forum.HandleFunc("/{slug}/users", h.forum.GetUsers).Methods(http.MethodGet)
	forum.HandleFunc("/{slug}/threads", h.forum.GetThreads).Methods(http.MethodGet)

	post := router.PathPrefix("/api/post").Subrouter()
	post.HandleFunc("/{id}/details", h.post.GetDetails).Methods(http.MethodGet)
	post.HandleFunc("/{id}/details", h.post.UpdateDetails).Methods(http.MethodPost)

	service := router.PathPrefix("/api/service").Subrouter()
	service.HandleFunc("/clear", h.service.ClearDb).Methods(http.MethodPost)
	service.HandleFunc("/status", h.service.StatusDb).Methods(http.MethodGet)

	thread := router.PathPrefix("/api/thread").Subrouter()
	thread.HandleFunc("/{slug_or_id}/create", h.post.CreatePosts).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/details", h.thread.GetDetails).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/details", h.thread.UpdateDetails).Methods(http.MethodPost)
	thread.HandleFunc("/{slug_or_id}/posts", h.thread.GetPosts).Methods(http.MethodGet)
	thread.HandleFunc("/{slug_or_id}/vote", h.thread.Vote).Methods(http.MethodPost)

	return router
}

func main() {
	logger.InitLogger()

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
			MaxConnections: 100,
			AfterConnect:   nil,
			AcquireTimeout: 0,
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
		user:    userHandler,
		forum:   forumHandler,
		post:    postHandler,
		service: serviceHandler,
		thread:  threadHandler,
	}

	router := newRouter(handlers)

	server := &http.Server{
		Handler: router,
		Addr:    ":5000",
	}

	logger.Start().Error(nil, errors.New("Server starting"))
	if err := server.ListenAndServe(); err != nil {
		logger.Start().Error(nil, err)
	}
}
