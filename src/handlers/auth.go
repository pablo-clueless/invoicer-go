package handlers

import (
	"invoicer-go/m/src/config"
	"invoicer-go/m/src/database"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/services"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		service: *services.NewAuthService(database.GetDatabase()),
	}
}

func (h *AuthHandler) SigninWithOauth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
		if err != nil {
			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
			return
		}
		response, err := h.service.SigninWithOauth(&user)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}
		lib.Success(ctx, "User signed in successfully", response)
	}
}

func (h *AuthHandler) GoogleCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (h *AuthHandler) SigninOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)

		err := h.service.SignOut(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "User signed out successfully", nil)
	}
}
