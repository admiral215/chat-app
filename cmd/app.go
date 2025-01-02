package cmd

import (
	"chat-app/config"
	_ "chat-app/docs"
	"chat-app/internal/delivery/http"
	"chat-app/internal/delivery/middleware"
	websocket2 "chat-app/internal/delivery/websocket"
	"chat-app/pkg/helpers"
	"chat-app/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	httpPkg "net/http"
)

type App struct {
	config     *config.Config
	jwtService jwt.JWTService
	hub        *websocket2.Hub

	authMiddleware *middleware.AuthMiddleware

	authHandler  *http.AuthHandler
	userHandler  *http.UserHandler
	groupHandler *http.GroupHandler
}

func NewApp(
	config *config.Config,
	jwtService jwt.JWTService,
	hub *websocket2.Hub,

	authMiddleware *middleware.AuthMiddleware,

	authHandler *http.AuthHandler,
	userHandler *http.UserHandler,
	groupHandler *http.GroupHandler,
) *App {
	return &App{
		config:     config,
		jwtService: jwtService,
		hub:        hub,

		authMiddleware: authMiddleware,

		userHandler:  userHandler,
		authHandler:  authHandler,
		groupHandler: groupHandler,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (app *App) Start() error {
	r := gin.Default()

	r.GET("/swagger/*any", gin.WrapF(httpSwagger.WrapHandler))

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", app.authHandler.Login)
			auth.POST("/register", app.userHandler.Register)
		}

		groups := api.Group("/groups", app.authMiddleware.Handle())
		{
			groups.POST("/", app.groupHandler.CreateGroup)
		}

		api.GET("/ws", app.authMiddleware.Handle(), func(c *gin.Context) {
			userId := helpers.GetUserFromContext(c)
			if userId == "" {
				c.JSON(500, gin.H{
					"error": "user id is empty",
				})
				return
			}

			serveWs(userId, app.hub, c.Writer, c.Request)
		})
	}

	go app.hub.Run()

	return r.Run(":" + app.config.App.Port)
}

// serveWs godoc
// @Summary WebSocket Connection (Can't to try)
// @Description Establish a WebSocket connection for real-time communication
// @Tags WebSocket
// @Param Authorization header string true "Bearer {token}"
// @Produce plain
// @Success 101 {string} string "Switching Protocols"
// @Router /api/ws [get]
func serveWs(userId string, hub *websocket2.Hub, w httpPkg.ResponseWriter, r *httpPkg.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &websocket2.Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userId,
	}

	client.Hub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()
}

// SendMessage godoc
// @Summary Send a message via WebSocket (Can't to try)
// @Description Send a JSON payload through the established WebSocket connection.
// @Description Example JSON payload for sending a message:
// @Description ```json
// @Description {
// @Description   "type": "private/group",
// @Description   "content": "test chat",
// @Description   "recipient_id": "group_id/user_id"
// @Description }
// @Tags WebSocket
// @Param Authorization header string true "Bearer {token}"
// @Router /api/ws/send [post]
func SendMessage(c *gin.Context) {
	// This is a placeholder to represent the payload schema and connection flow.
	c.JSON(httpPkg.StatusNotImplemented, gin.H{"message": "WebSocket message sending is done via the connection."})
}
