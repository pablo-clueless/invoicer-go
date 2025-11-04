package routes

import (
	"invoicer-go/m/src/handlers"

	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	invoices := router.Group("/invoices")
	handler := handlers.NewInvoiceHandler()

	invoices.POST("", handler.CreateInvoice())
	invoices.PUT("/:id", handler.UpdateInvoice())
	invoices.DELETE("/:id", handler.DeleteInvoice())
	invoices.GET("", handler.GetInvoices())
	invoices.GET("/:id", handler.GetInvoice())

	return invoices
}
