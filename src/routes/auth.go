package routes

import (
	"invoicer-go/m/src/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	auth := router.Group("/auth")
	handler := handlers.NewAuthHandler()

	auth.GET("/:provider", handler.SigninWithOauth())
	auth.GET("/:provider/callback", handler.GoogleCallback())
	auth.POST("/signout", handler.SigninOut())

	return auth
}
