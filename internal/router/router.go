package router

import (
	"DBProject/internal/handlers"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

func initRouter(handlers *handlers.Handlers) *router.Router {
	r := router.New()

	forumRoutes := r.Group(strings.Join([]string{RootRoute, ForumRoute}, ""))
	{
		forumRoutes.POST("/create", handlers.ForumHandler.CreateForum)
		forumRoutes.GET("/{slug}/details", handlers.ForumHandler.GetForum)
		forumRoutes.POST("/{slug}/create", handlers.ForumHandler.CreateThread)
		forumRoutes.GET("/{slug}/users", handlers.ForumHandler.GetForumUsers)
		forumRoutes.GET("/{slug}/{threads}", handlers.ForumHandler.GetForumThreads)
	}
	postRoutes := r.Group(strings.Join([]string{RootRoute, PostRoute}, ""))
	{
		postRoutes.GET("/{id}/details", handlers.PostHandler.GetPost)
		postRoutes.POST("/{id}/details", handlers.PostHandler.UpdatePost)
	}
	serviceRoutes := r.Group(strings.Join([]string{RootRoute, ServiceRoute}, ""))
	{
		serviceRoutes.POST("/clear", handlers.ServiceHandler.Clear)
		serviceRoutes.GET("/status", handlers.ServiceHandler.GetStatus)
	}
	threadRoutes := r.Group(strings.Join([]string{RootRoute, ThreadRoute}, ""))
	{
		threadRoutes.POST("/{slug_or_id}/create", handlers.ThreadHandler.CreatePosts)
		threadRoutes.GET("/{slug_or_id}/details", handlers.ThreadHandler.GetThread)
		threadRoutes.POST("/{slug_or_id}/details", handlers.ThreadHandler.UpdateThread)
		threadRoutes.GET("/{slug_or_id}/posts", handlers.ThreadHandler.GetThreadPosts)
		threadRoutes.POST("/{slug_or_id}/vote", handlers.ThreadHandler.Vote)
	}
	userRoutes := r.Group(strings.Join([]string{RootRoute, UserRoute}, ""))
	{
		userRoutes.POST("/{nickname}/create", handlers.UserHandler.CreateUser)
		userRoutes.GET("/{nickname}/profile", handlers.UserHandler.GetUser)
		userRoutes.POST("/{nickname}/profile", handlers.UserHandler.UpdateUser)
	}
	return r
}

func NewEngine(handlers *handlers.Handlers) *Engine {

	r := initRouter(handlers)
	handler := r.Handler

	engine := &Engine{
		router: r,
		server: &fasthttp.Server{
			Handler:            handler,
			ReadTimeout:        10 * time.Second,
			WriteTimeout:       10 * time.Second,
			MaxRequestBodySize: 1024 * 1024 * 1024,
		},
	}

	return engine
}

type Engine struct {
	server *fasthttp.Server
	router *router.Router
}

func (e *Engine) Start(port string) error {
	return e.server.ListenAndServe(port)
}

func (e *Engine) Stop() error {
	return e.server.Shutdown()
}
