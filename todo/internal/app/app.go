package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"todo/internal/adapter/controller/http/handler"
	"todo/internal/adapter/controller/http/middleware"
	"todo/internal/adapter/repository/postgres"
	sso "todo/internal/adapter/sso/grpc"
	"todo/internal/config"
	"todo/internal/domain/usecase"
)

type App struct {
	httpServer *http.Server
	log        *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) *App {
	// GRPC Client
	client, err := sso.New(
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		panic(err)
	}

	// Postgres
	conn, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresDB.Username,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName,
		),
	)
	if err != nil {
		panic(err)
	}

	// Repositories
	taskRep := postgres.NewTask(conn, log)
	listRep := postgres.NewList(conn, log)

	// Use cases
	listUseCase := usecase.NewList(listRep, log)
	taskUseCase := usecase.NewTask(taskRep, listUseCase, log)

	// Handlers
	authHandler := handler.NewAuthHandler(client, cfg.AppId)
	listHandler := handler.NewListHandler(listUseCase, log)
	taskHandler := handler.NewTaskHandler(taskUseCase, log)

	router := gin.New()
	auth := router.Group("/server")
	{
		auth.POST("/sign-up", authHandler.SignUp())
		auth.POST("/sign-in", authHandler.SignIn())
	}
	api := router.Group("/api/", middleware.AuthMiddleware([]byte(cfg.AppSecret), client))
	{
		lists := api.Group("/lists/")
		{
			lists.POST("/", listHandler.Create())
			lists.GET("/", listHandler.GetAll())
			lists.GET("/:id", listHandler.GetByID())
			lists.PUT("/:id", listHandler.UpdateByID())
			lists.DELETE("/:id", listHandler.DeleteByID())
			items := lists.Group(":id/tasks")
			{
				items.POST("/", taskHandler.Create())
				items.GET("/", taskHandler.GetByListID())
			}
		}
		items := api.Group("/tasks/")
		{
			items.GET("/:id", taskHandler.GetByID())
			items.PUT("/:id", taskHandler.UpdateByID())
			items.DELETE("/:id", taskHandler.DeleteByID())
		}
	}

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		httpServer: srv,
		log:        log,
	}
}

func (app *App) MustRun() error {
	err := app.run()
	return err
}

func (app *App) Stop(ctx context.Context) error {
	app.log.Info("Stopping HTTP server")
	return app.httpServer.Shutdown(ctx)
}

func (app *App) run() error {
	app.log.Info("Listening on ", app.httpServer.Addr)
	return app.httpServer.ListenAndServe()
}
