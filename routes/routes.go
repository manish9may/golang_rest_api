package routes

import (
	"golang_rest_api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", middleware.JWTAuthMiddleware(), fetchEvents)
	server.POST("/events", middleware.JWTAuthMiddleware(), createEvent)
	server.GET("/events/:id", middleware.JWTAuthMiddleware(), fetchEvent)
	server.PUT("/events/:id", middleware.JWTAuthMiddleware(), updateEvent)
	server.DELETE("/events/:id", middleware.JWTAuthMiddleware(), deleteEvent)
	server.POST("/signup", createUser)
	server.POST("/login", loginUser)
}
