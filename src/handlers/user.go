package handlers

import (
	"invoicer-go/m/src/config"
	"invoicer-go/m/src/database"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/lib"
	"invoicer-go/m/src/services"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		service: services.NewUserService(database.GetDatabase()),
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const (
			maxMemory = 10 << 20 // 10 MB
			bucket    = "companies"
		)

		contentType := ctx.ContentType()
		if !strings.Contains(contentType, "multipart/form-data") {
			lib.BadRequest(ctx, "Content type must be multipart/form-data", "400")
			return
		}

		if err := ctx.Request.ParseMultipartForm(maxMemory); err != nil {
			lib.BadRequest(ctx, "Failed to parse form data: "+err.Error(), "400")
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			lib.BadRequest(ctx, "Failed to get multipart form: "+err.Error(), "400")
			return
		}

		id := ctx.GetString(config.AppConfig.CurrentUserId)
		if id == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		var companyLogoPtr *string
		if len(form.File["companyLogo"]) > 0 {
			image := form.File["companyLogo"][0]
			if image != nil {
				companyLogo, uploadErr := lib.SingleImageUploader(image, bucket)
				if uploadErr != nil {
					lib.BadRequest(ctx, "Failed to upload company logo: "+uploadErr.Error(), "400")
					return
				}
				companyLogoPtr = &companyLogo
			}
		}

		payload := &dto.UpdateUserDto{
			Name:        getFormValue(form, "name"),
			Email:       getFormValue(form, "email"),
			Phone:       getFormValue(form, "phone"),
			RcNumber:    getFormValue(form, "rcNumber"),
			CompanyLogo: companyLogoPtr,
			CompanyName: getFormValue(form, "companyName"),
			Website:     getFormValue(form, "website"),
			TaxId:       getFormValue(form, "taxId"),
		}

		bankInfo := extractBankInformation(form)
		if bankInfo != nil {
			payload.BankInformation = bankInfo
		}

		if !hasUpdateFields(payload) {
			lib.BadRequest(ctx, "No valid fields to update", "400")
			return
		}

		user, err := h.service.UpdateUser(id, *payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Profile updated successfully", user)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		if id == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		confirm := ctx.Query("confirm")
		if strings.ToLower(confirm) != "true" {
			lib.BadRequest(ctx, "Deletion requires confirmation. Add ?confirm=true to proceed", "400")
			return
		}

		err := h.service.DeleteUser(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Account deleted successfully", nil)
	}
}

func (h *UserHandler) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.GetString(config.AppConfig.CurrentUserId)
		if id == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		user, err := h.service.GetUser(id)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Profile fetched successfully", user)
	}
}

func getFormValue(form *multipart.Form, key string) *string {
	if values, exists := form.Value[key]; exists && len(values) > 0 && strings.TrimSpace(values[0]) != "" {
		value := strings.TrimSpace(values[0])
		return &value
	}
	return nil
}

func extractBankInformation(form *multipart.Form) *dto.UpdateBankInformationDto {
	accountName := getFormValue(form, "bankInformation[accountName]")
	accountNumber := getFormValue(form, "bankInformation[accountNumber]")
	bankName := getFormValue(form, "bankInformation[bankName]")
	bankSwiftCode := getFormValue(form, "bankInformation[bankSwiftCode]")
	iban := getFormValue(form, "bankInformation[iban]")
	routingNumber := getFormValue(form, "bankInformation[routingNumber]")

	if accountName != nil || accountNumber != nil || bankName != nil ||
		bankSwiftCode != nil || iban != nil || routingNumber != nil {
		return &dto.UpdateBankInformationDto{
			AccountName:   accountName,
			AccountNumber: accountNumber,
			BankName:      bankName,
			BankSwiftCode: bankSwiftCode,
			Iban:          iban,
			RoutingNumber: routingNumber,
		}
	}

	return nil
}

func hasUpdateFields(payload *dto.UpdateUserDto) bool {
	return payload.Name != nil || payload.Email != nil || payload.Phone != nil ||
		payload.RcNumber != nil || payload.CompanyLogo != nil || payload.CompanyName != nil ||
		payload.Website != nil || payload.TaxId != nil || payload.BankInformation != nil
}

func (h *UserHandler) UpdateUserProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload dto.UpdateUserDto

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			lib.BadRequest(ctx, "Invalid request payload: "+err.Error(), "400")
			return
		}

		id := ctx.GetString(config.AppConfig.CurrentUserId)
		if id == "" {
			lib.Unauthorized(ctx, "User not authenticated")
			return
		}

		if !hasUpdateFields(&payload) {
			lib.BadRequest(ctx, "No valid fields to update", "400")
			return
		}

		user, err := h.service.UpdateUser(id, payload)
		if err != nil {
			lib.InternalServerError(ctx, err.Error())
			return
		}

		lib.Success(ctx, "Profile updated successfully", user)
	}
}
