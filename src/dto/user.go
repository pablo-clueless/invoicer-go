package dto

type UpdateUserDto struct {
	BankInformation *BankInformation `json:"bankInformation,omitempty"`
	CompanyLogo     *string          `json:"companyLogo,omitempty"`
	CompanyName     *string          `json:"companyName,omitempty"`
	Email           *string          `json:"email,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Phone           *string          `json:"phone,omitempty"`
	RcNumber        *string          `json:"rcNumber,omitempty"`
	TaxId           *string          `json:"taxId,omitempty"`
	Website         *string          `json:"website,omitempty"`
}

type BankInformation struct {
	AccountName   *string `json:"accountName,omitempty"`
	AccountNumber *string `json:"accountNumber,omitempty"`
	BankName      *string `json:"bankName,omitempty"`
	BankSwiftCode *string `json:"bankSwiftCode,omitempty"`
	Iban          *string `json:"iban,omitempty"`
	RoutingNumber *string `json:"routingNumber,omitempty"`
}
