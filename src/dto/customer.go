package dto

type CreateCustomerDto struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}
