package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"todo/internal/config"
	"todo/internal/storage/postgres"
	"todo/internal/transport/client/sso/grpc"
)

type Service struct {
	Task TaskStorage
	List ListStorage
}

type AppInfo struct {
	id     int
	secret string
}

type Handler struct {
	AppInfo
	log     *slog.Logger
	client  *grpc.Client
	service Service
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	client *grpc.Client,
	taskStorage *postgres.TaskStorage,
	listStorage *postgres.ListStorage,
) *Handler {
	return &Handler{
		AppInfo: AppInfo{
			cfg.AppId,
			cfg.AppSecret,
		},
		log:    log,
		client: client,
		service: Service{
			Task: taskStorage,
			List: listStorage,
		},
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListByID)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)
			items := lists.Group(":id/tasks")
			{
				items.POST("/", h.createTask)
				items.GET("/", h.getAllTasksByList)
			}
		}
		items := api.Group("/tasks")
		{
			items.GET("/:id", h.getTaskByID)
			items.PUT("/:id", h.updateTask)
			items.DELETE("/:id", h.deleteTask)
		}
	}
	return router
}
