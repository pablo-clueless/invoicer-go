package handlers

import (
	"invoicer-go/m/src/database"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/services"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service services.CustomerService
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		service: *services.NewCustomerService(database.GetDatabase()),
	}
}

func (h *CustomerHandler) CreateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.CreateCustomerDto

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		customer, err := h.service.CreateCustomer(payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customer created successfully", customer)
	}
}

func (h *CustomerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateCustomerDto
		id := ctx.Param("id")

		if err := ctx.ShouldBind(&payload); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		customer, err := h.service.UpdateCustomer(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customer updated succesfully", customer)
	}
}

func (h *CustomerHandler) DeleteCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		err := h.service.DeleteCustomer(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customer deleted successfully", nil)
	}
}

func (h *CustomerHandler) GetCustomers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params dto.CustomerPagination

		if err := ctx.ShouldBind(&params); err != nil {
			lib.BadRequest(ctx, err.Error(), "400")
			return
		}

		customers, err := h.service.GetCustomers(params)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customers fetched successfully", customers)
	}
}

func (h *CustomerHandler) GetCustomer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		customer, err := h.service.GetCustomer(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
		}
		lib.Success(ctx, "Customer fetched successfully", customer)
	}
}
