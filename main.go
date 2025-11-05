package main

import (
	"fmt"
	"invoicer-go/m/src/config"
	"invoicer-go/m/src/database"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/middlewares"
	"invoicer-go/m/src/routes"
	"invoicer-go/m/src/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitializeConfig()

	err := database.InitializeDatabase()
	defer database.CloseDatabase()
	if err != nil {
		log.Fatal("Database error:", err)
	}

	lib.InitialiseJWT(string(config.AppConfig.JWTSecret))

	if err := services.InitializeProvider(); err != nil {
		log.Fatal("OAuth provider initialization error:", err)
	}

	if config.AppConfig.IsDevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	app.Use(gin.Logger())
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "X-RateLimit-Limit", "X-RateLimit-Reset",
		},
		AllowMethods: []string{"DELETE", "GET", "POST", "PUT", "OPTIONS"},
		AllowOrigins: []string{"*", config.AppConfig.ClientUrl},
	}))
	app.Use(middlewares.ErrorHandlerMiddleware())
	app.Use(middlewares.AuthMiddleware())
	app.Use(lib.ErrorHandler())

	app.MaxMultipartMemory = 10 << 20 // 10MB

	hub := lib.NewHub()
	go hub.Run()

	prefix := config.AppConfig.Version
	router := app.Group(prefix)
	websocket := lib.NewWebSocketHandler(hub)

	router.GET("/ws", websocket.HandleWebSocket)
	router.GET("/health", func(ctx *gin.Context) {
		lib.Success(ctx, "Invoicer API is healthy", map[string]interface{}{
			"version": config.AppConfig.Version,
			"status":  http.StatusOK,
		})
	})

	routes.AuthRoutes(router)
	routes.CustomerRoutes(router)
	routes.InvoiceRoutes(router)
	routes.UserRoutes(router)

	app.NoRoute(lib.GlobalNotFound())

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.AppConfig.Port),
		Handler:        app,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("Server starting on port http://localhost:%s/%s", config.AppConfig.Port, config.AppConfig.Version)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
