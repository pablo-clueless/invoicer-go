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
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	BankName      string `json:"bankName"`
	BankSwiftCode string `json:"bankSwiftCode"`
	Iban          string `json:"iban"`
	RoutingNumber string `json:"routingNumber"`
}
