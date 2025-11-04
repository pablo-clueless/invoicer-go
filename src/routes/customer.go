package routes

import (
	"invoicer-go/m/src/handlers"

	"github.com/gin-gonic/gin"
)

func CustomerRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	customers := router.Group("/customers")
	handler := handlers.NewCustomerHandler()

	customers.POST("", handler.CreateCustomer())
	customers.PUT("/:id", handler.UpdateCustomer())
	customers.DELETE("/:id", handler.DeleteCustomer())
	customers.GET("", handler.GetCustomers())
	customers.GET("/:id", handler.GetCustomer())

	return customers
}
