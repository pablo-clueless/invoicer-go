package handlers

import (
	"invoicer-go/m/src/database"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/services"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	service services.InvoiceService
}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{
		service: *services.NewInvoiceService(database.GetDatabase()),
	}
}

func (h *InvoiceHandler) CreateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateInvoiceDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		invoice, err := h.service.CreateInvoice(payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Invoice created successfully", invoice)
	}
}

func (h *InvoiceHandler) UpdateInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateInvoiceDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		invoice, err := h.service.UpdateInvoice(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Invoice updated succesfully", invoice)
	}
}

func (h *InvoiceHandler) DeleteInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteInvoice(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customer deleted successfully", nil)
	}
}

func (h *InvoiceHandler) GetInvoices() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.InvoicePagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		invoices, err := h.service.GetInvoices(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Invoices fetched successfully", invoices)
	}
}

func (h *InvoiceHandler) GetInvoice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		invoice, err := h.service.GetInvoice(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Invoice fetched successfully", invoice)
	}
}
