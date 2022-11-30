package models

type Company struct {
	UUID              string
	Name              string
	AmountOfEmployees int64
	Registered        bool
	Type              string
	Description       string
}

type CompanyForUpdate struct {
	UUID              string
	Name              *string
	AmountOfEmployees *int64
	Registered        *bool
	Type              *string
	Description       *string
}
