package routes

import (
	"example.com/middlewares"
	"github.com/gin-gonic/gin"
)

// Prevent user from registering from his events .
// Check if user is registered already .
// Prevent user from deleting events with registration .
// Allow user to view all registrations .
// Standardize the error handling for db operations
// Standardize the responses

func RegisterRoutes(server *gin.Engine) {
	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	authenticated.POST("/events", createEvent)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)
	authenticated.GET("/events/:id/register", getRegistrations)
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", cancelRegistration)

	server.POST("/signup", signup)
	server.POST("/login", login)
}
