package routes

import (
	"invoicer-go/m/src/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	users := router.Group("/users")
	handler := handlers.NewUserHandler()

	users.PUT("", handler.UpdateUser())
	users.DELETE("", handler.DeleteUser())
	users.GET("", handler.GetUser())

	return users
}
