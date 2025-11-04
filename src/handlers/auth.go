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
		provider := ctx.Param("provider")

		if provider == "" {
			lib.BadRequest(ctx, "Provider is required", "400")
			return
		}

		q := ctx.Request.URL.Query()
		q.Add("provider", provider)
		ctx.Request.URL.RawQuery = q.Encode()

		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	}
}

func (h *AuthHandler) GoogleCallback() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.Param("provider")
		if provider == "" {
			lib.BadRequest(ctx, "Provider is required", "400")
			return
		}

		q := ctx.Request.URL.Query()
		q.Add("provider", provider)
		ctx.Request.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to complete OAuth authentication: "+err.Error())
			return
		}

		response, err := h.service.SigninWithOauth(&user)
		if err != nil {
			lib.InternalServerError(ctx, "Failed to sign in with OAuth: "+err.Error())
			return
		}
		lib.Success(ctx, "User signed in successfully", response)
	}
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
