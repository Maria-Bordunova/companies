package entity

type CompanyType string

const (
	corporations       CompanyType = "Corporations"
	nonProfit          CompanyType = "NonProfit"
	cooperative        CompanyType = "Cooperative"
	soleProprietorship CompanyType = "Sole Proprietorship"
)

type Company struct {
	UId         string
	Name        string
	Description *string
	Employees   int64
	Registered  bool
	Type        CompanyType
}

type CreateCompany Company

type UpdateCompany struct {
	Name        *string
	Description *string
	Employees   *int64
	Registered  *bool
	Type        *CompanyType
}
